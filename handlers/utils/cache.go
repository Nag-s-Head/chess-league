package utils

import (
	"fmt"
	"net/http"
)

const (
	AgeHour = 60 * 60
	AgeDay  = AgeHour * 24
	AgeWeek = AgeDay * 7
)

// This is used to give a page a cache age of 1 hour, use on statically rendered pages
func WithCacheControl(w http.ResponseWriter, age int) {
	w.Header().Set("cache-control", fmt.Sprintf("public, max-age=%d, s-max-age=%d", age, age))
}
