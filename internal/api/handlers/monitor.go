package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"queueit/internal/helper"
	"queueit/pkg/logger"
	"time"
)

var starttime = time.Now()

// healthResponse represents the response for the health check endpoint
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
}

// HandleHealth godoc
// @Summary Health check
// @Description Returns the server health status along with version and uptime
// @Tags Health
// @Produce json
// @Success 200 {object} healthResponse "Server is healthy"
// @Failure 500 {string} string "Health check failed"
// @Router /v1/health [get]
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	helper.SetJSONHeader(w)

	response := healthResponse{
		Status:  "UP",
		Version: "v1",
		Uptime:  fmt.Sprintf("%.f (second)", time.Since(starttime).Seconds()),
	}

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(response); err != nil {
		logger.Error("HandleHealth ~ JSON encoding failed: ")
		http.Error(w, "health check failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buffer.Bytes())
}
