package helper

import (
	"net/http"
	"os"
	"path/filepath"
	"queueit/internal/models"
)

// sets header content-type to JSON
//
// http writer to be sent as parameter
func SetJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// sets header content-type to JSON
//
// http writer to be sent as parameter
func SetTextHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
}

// fetch path for storing sqllite db related files
//
// LINUX: /home/user/.queueit/queueit.db
//
// WINDOWS: C:\Users\User\AppData\Roaming\queueit\queueit.db
func GetAppDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if os.PathSeparator == '\\' { // Windows
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "queueit"), nil
		}
	}
	return filepath.Join(home, ".queueit"), nil
}

func IsValidStatus(s int) bool {
	return models.ValidStatuses[s]
}

func IsValidPriority(s int) bool {
	return models.ValidPriorities[s]
}
