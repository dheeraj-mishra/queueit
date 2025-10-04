package api

import (
	"net/http"
	"queueit/internal/api/handlers"

	"github.com/gorilla/mux"
)

type API struct {
	router *mux.Router
}

// creates a new router instance usingn gorilla mux lib
func NewRouter() *API {
	mr := mux.NewRouter()
	mr.HandleFunc("/v1/health", handlers.HandleHealth).Methods("GET")
	mr.HandleFunc("/v1/tasks", handlers.GetAllTasks).Methods("GET")
	// mr.HandleFunc("/api/v1/tasks/{id}", handlers.GetTaskByID).Methods("GET")
	mr.HandleFunc("/v1/tasks", handlers.CreateTask).Methods("POST")
	// mr.HandleFunc("/api/v1/tasks/{id}", handlers.UpdateTask).Methods("PUT", "PATCH")
	// mr.HandleFunc("/api/v1/tasks/{id}", handlers.DeleteTask).Methods("DELETE")
	// mr.HandleFunc("/api/v1/tasks/{id}/deadline", handlers.UpdateDeadline).Methods("PATCH")
	// mr.HandleFunc("/api/v1/tasks/{id}/status", handlers.UpdateStatus).Methods("PATCH")

	return &API{
		router: mr,
	}
}

func (api API) StartServer() error {
	return http.ListenAndServe(":6666", api.router)
}
