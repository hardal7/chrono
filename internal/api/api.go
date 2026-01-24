package api

import (
	"net/http"

	"github.com/hardal7/chrono/internal/config"
	"github.com/hardal7/chrono/internal/handler/topic"
	"github.com/hardal7/chrono/internal/handler/user"
	"github.com/hardal7/chrono/internal/middleware"
	logger "github.com/hardal7/chrono/internal/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunAPIServer() {
	admin := http.NewServeMux()
	admin.Handle("/metrics", promhttp.Handler())

	root := http.NewServeMux()

	public := http.NewServeMux()
	public.HandleFunc("POST /register", CreateRequest(user.Register, "register user"))
	public.HandleFunc("POST /login", CreateRequest(user.Login, "log user in"))
	root.Handle("/register", public)
	root.Handle("/login", public)

	protected := http.NewServeMux()
	protected.HandleFunc("POST /account", CreateRequest(user.EditAccount, "edit user account"))
	protected.HandleFunc("POST /topic/create", CreateRequest(topic.Create, "create topic"))
	protected.HandleFunc("POST /topic/edit", CreateRequest(topic.Edit, "edit topic"))
	protected.HandleFunc("POST /topic/track", CreateRequest(topic.Track, "track topic"))
	protected.HandleFunc("GET /topic/events", CreateRequest(topic.GetEvents, "get topic events"))
	root.Handle("/", middleware.Authenticate(protected))

	server := http.Server{
		Addr:    ":" + config.App.Port,
		Handler: middleware.LogRequest(root),
	}

	adminServer := http.Server{
		Addr:    ":" + config.App.AdminPort,
		Handler: middleware.LogRequest(admin),
	}

	go func() {
		logger.Info("Starting admin server on port: " + config.App.AdminPort)
		err := adminServer.ListenAndServe()
		if err != nil {
			logger.Error("Failed to start admin server")
			logger.Debug(err.Error())
		}
	}()

	logger.Info("Starting server on port: " + config.App.Port)
	err := server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server")
		logger.Debug(err.Error())
	}
}
