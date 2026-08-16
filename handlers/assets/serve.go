package assets

import "net/http"

func ServeAsset(bytes []byte, contentType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", contentType)
		w.Header().Set("cache-control", "public, max-age=604800, s-max-age=604800")
		w.Write(bytes)
	}
}
