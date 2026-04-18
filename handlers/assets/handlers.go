package assets

import (
	"net/http"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

func serveAsset(data []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
		w.Header().Set("Content-Type", "image/jpeg")
		utils.WithCacheControl(w, utils.AgeWeek)
	}
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /assets/wide_shot.jpg", serveAsset(wideShot))
}
