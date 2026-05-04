package model

import (
	"net"
	"net/http"
	"strings"
)

// GetRemoteAddr returns the remote IP address of the request.
// It prioritizes the X-Forwarded-For header if present, otherwise it uses r.RemoteAddr.
// If r.RemoteAddr contains a port, it is stripped.
func GetRemoteAddr(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can be a comma-separated list. The first one is the client.
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
