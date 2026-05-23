package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/domains/health"
	"github.com/hardal7/chrono/internal/domains/topic"
	"github.com/hardal7/chrono/internal/domains/topicevent"
	"github.com/hardal7/chrono/internal/domains/user"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/handler"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve() {
	adminRouter := chi.NewRouter()
	adminRouter.Handle("/metrics", promhttp.Handler())

	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.LogRequest)

	// Public Routes
	mainRouter.Group(func(r chi.Router) {
		r.Post("/register", handler.Create(user.Register))
		r.Post("/login", handler.Create(user.Login))
		r.Get("/health", health.Ping)
	})

	// Protected Routes
	mainRouter.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Route("/user", user.Routes)
		r.Route("/topic", topic.Routes)
		r.Route("/topic-event", topicevent.Routes)
	})

	go runServer("main", config.App.Port, mainRouter)
	runServer("admin", config.App.AdminPort, adminRouter)
}

func runServer(name, port string, router *chi.Mux) {
	logger.Info("Starting " + name + " server on port: " + port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		logger.Fatal("Failed to start "+name+" server", err)
	}
}
