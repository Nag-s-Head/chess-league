package security

import (
	"crypto/rand"
	"strings"
)

func NewSessionkey() string {
	var sb strings.Builder
	for range 5 {
		sb.WriteString(rand.Text())
	}
	return sb.String()
}
