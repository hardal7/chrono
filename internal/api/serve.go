package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve() {
	InitValidator()

	adminRouter := chi.NewRouter()
	adminRouter.Handle("/metrics", promhttp.Handler())

	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.LogRequest)

	// Public routes
	mainRouter.Group(func(r chi.Router) {
		r.Post("/register", RegisterUserHandler)
		r.Post("/login", LoginUserHandler)
		r.Get("/health", PingHandler)
		r.Get("/avatar/{id}", GetUserAvatarHandler)
	})

	// Authenticated routes
	mainRouter.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Route("/user", UserRoute)
		r.Route("/topic", TopicRoute)
		r.Route("/topic-event", TopicEventRoute)
		r.Route("/session", SessionRoute)
		r.Route("/friend", FriendRoute)
	})

	go runServer("main", config.App.Port, mainRouter)
	runServer("admin", config.App.AdminPort, adminRouter)
}

func runServer(name, port string, router *chi.Mux) {
	logger.Info("Started HTTP server", "type", name, "port", ":"+port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		logger.Fatal("Fatal error on server", err, "type", name)
	}
}

var validate *validator.Validate

func InitValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
