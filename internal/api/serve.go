package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/hardal7/chrono/internal/auth"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const siteDir = "./static/site/"

func Serve(ctx context.Context) error {
	InitValidator()

	otelShutdown, err := telemetry.InitOTel(ctx)
	if err != nil {
		return err
	}
	defer otelShutdown(ctx)

	mainRouter := chi.NewRouter()
	mainRouter.Use(otelhttp.NewMiddleware("api-server"))
	mainRouter.Use(middleware.LogRequest)

	mainRouter.Route("/api", func(r chi.Router) {
		r.Group(publicRoutes)

		r.Group(func(r chi.Router) {
			r.Use(auth.Authenticate)
			r.Use(middleware.Activity)

			r.Route("/user", UserRoute)
			r.Route("/location", LocationRoute)
			r.Route("/topic", TopicRoute)
			r.Route("/topic-event", TopicEventRoute)
			r.Route("/session", SessionRoute)
			r.Route("/friend", FriendRoute)
		})
	})

	siteRoutes := []string{"privacy", "terms", "report", "feature"}
	for _, route := range siteRoutes {
		mainRouter.Get("/"+route, serveHTML(route+".html"))
	}

	siteServer := http.FileServer(http.Dir(siteDir))
	mainRouter.Handle("/*", siteServer)

	go runServer(ctx, "main", config.App.Port, mainRouter)

	return err
}

func runServer(ctx context.Context, name, port string, handler http.Handler) {
	server := &http.Server{
		ReadHeaderTimeout: time.Second,
		Addr:              ":" + port,
		Handler:           handler,
	}

	logger.With("type", name, "port", port).Info("Started HTTP server")

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("Error listening and serving HTTP server", err)
	}

	logger.With("type", name).Info("Shutting down HTTP server")
	err = server.Shutdown(ctx)
	if err == nil {
		logger.With("type", name).Info("Shut down HTTP server")
	} else {
		logger.Err(err).
			With("type", name).
			Error("shut down HTTP server")
	}
}

func serveHTML(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, siteDir+path)
	}
}

var validate *validator.Validate

func InitValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
