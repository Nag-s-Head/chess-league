package security

import "crypto/rand"

func NewSessionkey() string {
	return rand.Text()
}
