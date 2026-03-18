package model_test

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewAdminUser(t *testing.T) {
	name := uuid.New().String()
	oauthId := uuid.New().String()
	lastIp := "127.0.0.1"
	lastUserAgent := uuid.New().String()

	user := model.NewAdminUser(name, oauthId, lastIp, lastUserAgent)
	require.NotEmpty(t, user)
	require.Equal(t, name, user.Name)
	require.Equal(t, oauthId, user.OauthId)
	require.Equal(t, lastIp, user.LastIp)
	require.Equal(t, lastUserAgent, user.LastUserAgent)

	require.NotEmpty(t, user.Id)
	require.NotEmpty(t, user.SessionKey)
	require.True(t, len(user.SessionKey) > 12)

	require.NotEmpty(t, user.LastLogin)
	require.NotEmpty(t, user.Created)
}

func TestAdminLoginNewUser(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := "New Admin"
	oauthId := "github-" + uuid.New().String()
	ip := "192.168.1.1"
	ua := "Mozilla/5.0"

	user, err := model.AdminLogin(db, name, oauthId, ip, ua)
	require.NoError(t, err)
	require.Equal(t, name, user.Name)
	require.Equal(t, oauthId, user.OauthId)
	require.Equal(t, ip, user.LastIp)
	require.Equal(t, ua, user.LastUserAgent)
	require.NotEmpty(t, user.Id)
	require.NotEmpty(t, user.SessionKey)
	require.NotEmpty(t, user.LastLogin)
	require.NotEmpty(t, user.Created)

	// Verify database content
	var dbUser model.AdminUser
	err = db.GetSqlxDb().Get(&dbUser, "SELECT id, name, oauth_id, created, COALESCE(session_key, '') as session_key, last_login, last_ip, last_user_agent FROM admin_users WHERE oauth_id = $1", oauthId)
	require.NoError(t, err)
	require.Equal(t, user.Id, dbUser.Id)
	require.Equal(t, name, dbUser.Name)
	require.Equal(t, ip, dbUser.LastIp)
	require.Equal(t, ua, dbUser.LastUserAgent)
}

func TestAdminLoginExistingUser(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	oauthId := "github-" + uuid.New().String()

	// First login
	user1, err := model.AdminLogin(db, "First Name", oauthId, "1.1.1.1", "UA1")
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond) // Ensure time moves forward for LastLogin comparison

	// Second login with updated details
	newName := "Second Name"
	newIp := "2.2.2.2"
	newUA := "UA2"

	user2, err := model.AdminLogin(db, newName, oauthId, newIp, newUA)
	require.NoError(t, err)
	require.Equal(t, user1.Id, user2.Id)
	require.Equal(t, newName, user2.Name)
	require.Equal(t, oauthId, user2.OauthId)
	require.Equal(t, newIp, user2.LastIp)
	require.Equal(t, newUA, user2.LastUserAgent)
	require.NotEqual(t, user1.SessionKey, user2.SessionKey)
	require.True(t, user2.LastLogin.After(user1.LastLogin))

	// Verify database content
	var dbUser model.AdminUser
	err = db.GetSqlxDb().Get(&dbUser, "SELECT id, name, oauth_id, created, COALESCE(session_key, '') as session_key, last_login, last_ip, last_user_agent FROM admin_users WHERE oauth_id = $1", oauthId)
	require.NoError(t, err)
	require.Equal(t, newName, dbUser.Name)
	require.Equal(t, newIp, dbUser.LastIp)
	require.Equal(t, newUA, dbUser.LastUserAgent)
	require.Equal(t, user2.SessionKey, dbUser.SessionKey)
}

func TestAdminGetFromSessionKey(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := "Session User"
	oauthId := "github-" + uuid.New().String()
	user, err := model.AdminLogin(db, name, oauthId, "1.1.1.1", "UA")
	require.NoError(t, err)

	retrieved, err := model.AdminGetFromSessionKey(db, user.SessionKey)
	require.NoError(t, err)
	require.Equal(t, user.Id, retrieved.Id)
	require.Equal(t, user.Name, retrieved.Name)
}

func TestAdminGetFromSessionKeyFails(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	t.Run("Empty key", func(t *testing.T) {
		_, err := model.AdminGetFromSessionKey(db, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("Non-existent key", func(t *testing.T) {
		_, err := model.AdminGetFromSessionKey(db, "this-key-does-not-exist")
		require.Error(t, err)
	})
}

func TestAdminGetFromSessionKeyExpired(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := "Expired User"
	oauthId := "github-" + uuid.New().String()
	user, err := model.AdminLogin(db, name, oauthId, "1.1.1.1", "UA")
	require.NoError(t, err)

	// Manually expire the session in the database
	pastLogin := time.Now().Add(-model.MaxSessionKeyAge - time.Minute)
	_, err = db.GetSqlxDb().Exec("UPDATE admin_users SET last_login = $1 WHERE id = $2", pastLogin, user.Id)
	require.NoError(t, err)

	_, err = model.AdminGetFromSessionKey(db, user.SessionKey)
	require.Error(t, err)
	require.Contains(t, err.Error(), "User has been logged in for too long")

	// Verify session key was cleared
	var dbUser model.AdminUser
	err = db.GetSqlxDb().Get(&dbUser, "SELECT id, name, oauth_id, created, COALESCE(session_key, '') as session_key, last_login, last_ip, last_user_agent FROM admin_users WHERE id = $1", user.Id)
	require.NoError(t, err)
	require.Empty(t, dbUser.SessionKey)
}

func TestAdminLogout(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	oauthId := "github-" + uuid.New().String()
	user, err := model.AdminLogin(db, "Logout User", oauthId, "1.1.1.1", "UA")
	require.NoError(t, err)
	require.NotEmpty(t, user.SessionKey)

	err = model.AdminLogout(db, user.Id)
	require.NoError(t, err)

	// Verify session key is gone
	var dbUser model.AdminUser
	err = db.GetSqlxDb().Get(&dbUser, "SELECT id, name, oauth_id, created, COALESCE(session_key, '') as session_key, last_login, last_ip, last_user_agent FROM admin_users WHERE id = $1", user.Id)
	require.NoError(t, err)
	require.Empty(t, dbUser.SessionKey)

	// Verify AdminGetFromSessionKey fails
	_, err = model.AdminGetFromSessionKey(db, user.SessionKey)
	require.Error(t, err)
}

func TestAdminLogoutNonExistentUser(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	// Should not error even if user doesn't exist
	err := model.AdminLogout(db, uuid.New())
	require.NoError(t, err)
}

func TestGetAdminUsers(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	const numberOfUsers = 100

	t.Log("Inserting lots of users")
	for range numberOfUsers {
		name := rand.Text()
		oauthId := name + "-id"
		user := model.NewAdminUser(name, oauthId, rand.Text(), rand.Text())

		tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
		require.NoError(t, err)
		defer tx.Rollback()

		require.NoError(t, model.AddAdminUser(tx, *user))
		require.NoError(t, tx.Commit())
	}

	t.Log("Actually doing the test")
	users, err := model.GetAdminUsers(db)
	require.NoError(t, err)
	require.True(t, len(users) >= numberOfUsers)

	for _, user := range users {
		require.NotEmpty(t, user)
		require.Empty(t, user.SessionKey)

		require.NotEmpty(t, user.Id)
		require.NotEmpty(t, user.Name)
		require.NotEmpty(t, user.OauthId)
		require.NotEmpty(t, user.Created)
		require.NotEmpty(t, user.LastIp)
		require.NotEmpty(t, user.LastLogin)
		require.NotEmpty(t, user.LastUserAgent)
	}
}
