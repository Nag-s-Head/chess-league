package utils

import (
	"net/http"
)

// This is used to give a page a cache age of 1 hour, use on statically rendered pages
func WithCacheControl(w http.ResponseWriter) {
	w.Header().Set("cache-control", "public, max-age=3600, s-max-age=3600")
}
