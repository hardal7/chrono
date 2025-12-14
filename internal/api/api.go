package api

import (
	"encoding/json"
	"net/http"

	"github.com/hardal7/chrono/internal/config"
	"github.com/hardal7/chrono/internal/handler/user"
	"github.com/hardal7/chrono/internal/middleware"
	logger "github.com/hardal7/chrono/internal/util"
)

func RunAPIServer() {
	root := http.NewServeMux()

	public := http.NewServeMux()
	public.HandleFunc("POST /register", CreateRequest(user.Register, "register user"))
	public.HandleFunc("POST /login", CreateRequest(user.Login, "log user in"))
	root.Handle("/register", public)
	root.Handle("/login", public)

	protected := http.NewServeMux()
	protected.HandleFunc("POST /account", CreateRequest(user.EditAccount, "edit user account"))
	// protected.HandleFunc("POST /session", CreateRequest(session.CreateSession, "create session"))
	root.Handle("/", middleware.Authenticate(protected))

	logger.Info("Starting server on port: " + config.App.Port)
	server := http.Server{
		Addr:    ":" + config.App.Port,
		Handler: middleware.LogRequest(root),
	}
	err := server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server")
		logger.Debug(err.Error())
	}
}

func CreateRequest[v any](f func(http.ResponseWriter, *http.Request, v), operation string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req v
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Info("Failed to " + operation)
			logger.Info("Failed to decode JSON")
			logger.Debug(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			f(w, r, req)
		}
	})
}
