package handlers

import (
	_ "embed"
	"net/http"
)

//go:embed webview/index.html
var indexHTML []byte

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(indexHTML)
}
