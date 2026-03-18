package model

import (
	"time"

	"github.com/google/uuid"
)

type AdminUser struct {
	Id            uuid.UUID `db:"id"`
	Name          string    `db:"name"`
	OauthId       string    `db:"oauth_id"` // This is the Github login
	Created       time.Time `db:"created"`
	SessionKey    string    `db:"session_key"`
	LastLogin     time.Time `db:"last_login"` // Used to know if the session key is expired
	LastIp        string    `db:"last_ip"`
	LastUserAgent string    `db:"last_user_agent"`
}
