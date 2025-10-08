package main

import (
	"fmt"
	"os"
	"queueit/internal/api"
	"queueit/internal/config"
	"queueit/internal/db"
	"queueit/pkg/logger"
	"time"

	webview "github.com/webview/webview_go"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		logger.Fatal(err)
	}

	if err := db.InitDB(); err != nil {
		logger.Fatal(err)
	}

	router := api.NewRouter()
	go func() {
		if err := router.StartServer(); err != nil {
			logger.Fatal(err)
		}
	}()

	// Wait a moment to ensure server starts
	time.Sleep(500 * time.Millisecond)

	// Start WebView window
	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("queueit")
	w.SetSize(1200, 800, webview.Hint(0))

	w.Navigate(fmt.Sprintf("http://%s:%s", os.Getenv("SERVER_IP"), os.Getenv("SERVER_PORT")))
	w.Run()
}
