package api

import (
	"net/http"

	"github.com/hardal7/study/internal/config"
	"github.com/hardal7/study/internal/middleware"
	logger "github.com/hardal7/study/internal/util"
)

func RunAPIServer() {
	router := http.NewServeMux()
	router.HandleFunc("POST /register", CreateRegisterRequest)
	router.HandleFunc("POST /login", CreateLoginRequest)
	router.HandleFunc("POST /account", CreateEditAccountRequest)

	logger.Info("Starting server on port: " + config.App.Port)
	server := http.Server{
		Addr:    ":" + config.App.Port,
		Handler: middleware.LogRequest(router),
	}
	err := server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server")
		logger.Debug(err.Error())
	}
}
