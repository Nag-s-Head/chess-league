package model

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/security"
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

func NewAdminUser(name, oauthId, lastIp, LastUserAgent string) AdminUser {
	return AdminUser{
		Id:            uuid.New(),
		Name:          name,
		OauthId:       oauthId,
		LastIp:        lastIp,
		LastUserAgent: LastUserAgent,
		LastLogin:     time.Now(),
		Created:       time.Now(),
		SessionKey:    security.NewSessionkey(),
	}
}

func AdminLogin(db *db.Db, name, oauthId, lastIp, LastUserAgent string) (AdminUser, error) {
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	if err != nil {
		return AdminUser{}, errors.Join(errors.New("Could not begin transaction"), err)
	}

	defer tx.Rollback()

	var user AdminUser
	err = tx.Get(&user, "SELECT * FROM admin_users WHERE oauth_id = $1;", oauthId)

	if err == nil {
		// Update with new session information
		user.Name = name
		user.SessionKey = security.NewSessionkey()
		user.LastLogin = time.Now()
		user.LastIp = lastIp
		user.LastUserAgent = LastUserAgent

		_, err = tx.NamedExec("UPDATE admin_users SET last_login=:last_login, session_key=:session_key, name=:name, last_ip=:last_ip, last_user_agent=:last_user_agent WHERE id=:id;", user)
		if err != nil {
			return AdminUser{}, errors.Join(errors.New("Could not update admin user"), err)
		}

		slog.Info("Admin user has logged in", "id", user.Id, "oauthId", user.OauthId, "name", user.Name)
	} else if errors.Is(err, sql.ErrNoRows) {
		user = NewAdminUser(name, oauthId, lastIp, LastUserAgent)
		_, err := tx.NamedExec(`
			INSERT INTO admin_users(id, name, oauth_id, created, session_key, last_login, last_ip, last_user_agent)
			VALUES (:id, :name, :oauth_id, :created, :session_key, :last_login, :last_ip, :last_user_agent);
			`, user)
		if err != nil {
			return AdminUser{}, errors.Join(errors.New("Could not create new admin user"), err)
		}

		slog.Info("Created a new admin user", "id", user.Id, "oauthId", user.OauthId, "name", user.Name)
	} else {
		return AdminUser{}, errors.Join(errors.New("Could not get the admin user"), err)
	}

	err = tx.Commit()
	if err != nil {
		return AdminUser{}, errors.Join(errors.New("Cold not commit transaction"), err)
	}

	return user, nil
}
