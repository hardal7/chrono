package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/config"
	"github.com/hardal7/chrono/internal/handler/health"
	"github.com/hardal7/chrono/internal/handler/topic"
	"github.com/hardal7/chrono/internal/handler/user"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunAPIServer() {
	adminRouter := chi.NewRouter()
	adminRouter.Handle("/metrics", promhttp.Handler())

	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.LogRequest)

	mainRouter.Group(func(r chi.Router) {
		r.Post("/register", CreateRequest(user.Register, "register user"))
		r.Post("/login", CreateRequest(user.Login, "log user in"))
		r.Get("/health", http.HandlerFunc(health.Ping))
	})

	mainRouter.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate)

		r.Post("/account", CreateRequest(user.EditAccount, "edit user account"))

		r.Route("/topics", func(r chi.Router) {
			r.Post("/create", CreateRequest(topic.Create, "create topic"))
			r.Post("/edit", CreateRequest(topic.Edit, "edit topic"))
			r.Post("/track", CreateRequest(topic.Track, "track topic"))
			r.Get("/events", CreateRequest(topic.GetEvents, "get topic events"))
		})
	})

	go runServer("main", config.App.Port, mainRouter)
	runServer("admin", config.App.AdminPort, adminRouter)
}

func runServer(name, port string, router *chi.Mux) {
	logger.Info("Starting " + name + " server on port: " + port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		logger.Error("Failed to start " + name + " server")
		logger.Debug(err.Error())
	}
}
