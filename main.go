package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const addr = "0.0.0.0:8080"

func main() {
	slog.Info("Starting Nag's Knights chess league application...", "addr", addr)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("alive and well at %s", time.Now().UTC()))
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("Could not start", "err", err, "addr", addr)
	}
}
