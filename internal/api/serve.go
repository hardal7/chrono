package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/hardal7/chrono/internal/auth"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(ctx context.Context) {
	InitValidator()

	adminRouter := chi.NewRouter()
	adminRouter.Handle("/metrics", promhttp.Handler())

	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.LogRequest)

	mainRouter.Route("/api", func(r chi.Router) {
		r.Group(publicRoutes)

		r.Group(func(r chi.Router) {
			r.Use(auth.Authenticate)
			r.Route("/user", UserRoute)
			r.Route("/location", LocationRoute)
			r.Route("/topic", TopicRoute)
			r.Route("/topic-event", TopicEventRoute)
			r.Route("/session", SessionRoute)
			r.Route("/friend", FriendRoute)
		})
	})

	mainRouter.Get("/privacy", serveHTML("privacy.html"))
	mainRouter.Get("/terms", serveHTML("terms.html"))
	mainRouter.Get("/report", serveHTML("report.html"))
	mainRouter.Get("/feature", serveHTML("feature.html"))

	siteServer := http.FileServer(http.Dir("./static"))
	mainRouter.Handle("/*", siteServer)

	go runServer(ctx, "main", config.App.Port, mainRouter)
	runServer(ctx, "admin", config.App.AdminPort, adminRouter)
}

func runServer(ctx context.Context, name, port string, router *chi.Mux) {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		logger.Info("Started HTTP server", "type", name, "port", ":"+port)

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Fatal errror on HTTP server", "type", name, "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down HTTP server", "type", name)
	err := server.Shutdown(ctx)
	if err == nil {
		logger.Info("Shut down HTTP server", "type", name)
	} else {
		logger.Error("Failed to shut down HTTP server", "type", name, "error", err)
	}
}

func serveHTML(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/"+path)
	}
}

var validate *validator.Validate

func InitValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
