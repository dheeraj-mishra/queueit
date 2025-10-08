package api

import (
	"fmt"
	"net/http"
	"os"
	"queueit/internal/api/handlers"
	"queueit/internal/api/middleware"
	"queueit/pkg/logger"

	"github.com/gorilla/mux"
)

type API struct {
	router *mux.Router
}

// creates a new router instance usingn gorilla mux lib
func NewRouter() *API {
	mr := mux.NewRouter()

	// middleware implementations:
	mr.Use(middleware.CORSMiddleware)
	mr.Use(middleware.LoggingMiddleware)

	mr.HandleFunc("/v1/health", handlers.HandleHealth).Methods("GET", "OPTIONS")
	mr.HandleFunc("/v1/tasks", handlers.GetAllTasks).Methods("GET", "OPTIONS")
	mr.HandleFunc("/v1/tasks/{id}", handlers.GetTaskByID).Methods("GET")
	mr.HandleFunc("/v1/tasks", handlers.CreateTask).Methods("POST", "OPTIONS")
	mr.HandleFunc("/v1/tasks/{id}", handlers.UpdateTask).Methods("PUT", "PATCH")
	mr.HandleFunc("/v1/tasks/{id}", handlers.DeleteTask).Methods("DELETE")
	mr.HandleFunc("/", handlers.Home)

	logger.Info("router created")
	return &API{
		router: mr,
	}
}

func (api API) StartServer() error {
	url_base := fmt.Sprintf("%s:%s", os.Getenv("SERVER_IP"), os.Getenv("SERVER_PORT"))

	logger.Info("router started, ready to accept requests")
	logger.Info("router ip:port", url_base)

	return http.ListenAndServe(url_base, api.router)
}
