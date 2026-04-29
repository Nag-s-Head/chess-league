package model_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestInsertAuditLog(t *testing.T) {
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	admin := model.NewAdminUser("bob", uuid.New().String(), "uwu", "uwu")
	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	name := uuid.New().String()
	desc := uuid.New().String()

	require.NoError(t, model.InsertAdminUser(tx, *admin))

	auditLog := model.NewAuditLog(admin.Id, name, desc)
	require.NotEmpty(t, auditLog.Id)
	require.NotEmpty(t, auditLog.Created)
	require.Equal(t, auditLog.OperationName, name)
	require.Equal(t, auditLog.OperationDescription, desc)

	require.NoError(t, model.InsertAuditLog(tx, auditLog))
	require.NoError(t, tx.Commit())
}

func TestInsertAuditLogPlayerAffected(t *testing.T) {
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	player := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayerTx(tx, player))

	admin := model.NewAdminUser("bob", uuid.New().String(), "uwu", "uwu")
	require.NoError(t, model.InsertAdminUser(tx, *admin))

	name := uuid.New().String()
	desc := uuid.New().String()

	auditLog := model.NewAuditLog(admin.Id, name, desc)
	require.NotEmpty(t, auditLog.Id)
	require.NotEmpty(t, auditLog.Created)
	require.Equal(t, auditLog.OperationName, name)
	require.Equal(t, auditLog.OperationDescription, desc)

	require.NoError(t, model.InsertAuditLog(tx, auditLog))
	require.NoError(t, model.InsertAuditLogPlayerAffected(tx, model.NewAuditLogPlayerAffected(auditLog.Id, player.Id, 123)))
	require.NoError(t, tx.Commit())
}

func TestGetAuditLog(t *testing.T) {
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	player := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayerTx(tx, player))

	admin := model.NewAdminUser("bob", uuid.New().String(), "uwu", "uwu")
	require.NoError(t, model.InsertAdminUser(tx, *admin))

	name := uuid.New().String()
	desc := uuid.New().String()

	auditLog := model.NewAuditLog(admin.Id, name, desc)
	require.NotEmpty(t, auditLog.Id)
	require.NotEmpty(t, auditLog.Created)
	require.Equal(t, auditLog.OperationName, name)
	require.Equal(t, auditLog.OperationDescription, desc)

	require.NoError(t, model.InsertAuditLog(tx, auditLog))
	require.NoError(t, model.InsertAuditLogPlayerAffected(tx, model.NewAuditLogPlayerAffected(auditLog.Id, player.Id, 123)))

	details, err := model.GetAuditLog(tx, auditLog.Id)
	require.NoError(t, err)
	require.Equal(t, auditLog.Id, details.Id)
	require.Equal(t, name, details.OperationName)
	require.Equal(t, desc, details.OperationDescription)
	require.Len(t, details.Players, 1)
	require.Equal(t, player.Id, details.Players[0].PlayerId)
	require.NoError(t, tx.Commit())
}

func TestGetAuditLogsUiFriendly(t *testing.T) {
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	admin := model.NewAdminUser("bob", uuid.New().String(), "uwu", "uwu")
	require.NoError(t, model.InsertAdminUser(tx, *admin))

	name := uuid.New().String()
	desc := uuid.New().String()

	auditLog := model.NewAuditLog(admin.Id, name, desc)
	require.NoError(t, model.InsertAuditLog(tx, auditLog))
	require.NoError(t, tx.Commit())

	auditLogs, err := model.GetAuditLogsUiFriendly(db)
	require.NoError(t, err)
	require.True(t, len(auditLogs) > 1)

	for _, log := range auditLogs {
		require.NotEmpty(t, log)
		require.NotEmpty(t, log.AdminName)
	}
}

func TestGetAuditLogsUiFriendlyByPlayer(t *testing.T) {
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	player := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayerTx(tx, player))

	admin := model.NewAdminUser("bob", uuid.New().String(), "uwu", "uwu")
	require.NoError(t, model.InsertAdminUser(tx, *admin))

	name := uuid.New().String()
	desc := uuid.New().String()

	auditLog := model.NewAuditLog(admin.Id, name, desc)
	require.NoError(t, model.InsertAuditLog(tx, auditLog))
	require.NoError(t, model.InsertAuditLogPlayerAffected(tx, model.NewAuditLogPlayerAffected(auditLog.Id, player.Id, 123)))
	require.NoError(t, tx.Commit())

	auditLogs, err := model.GetAuditLogsUiFriendlyForPlayer(db, player.Id)
	require.NoError(t, err)
	require.Len(t, auditLogs, 1)

	require.NotEmpty(t, auditLogs[0])
	require.NotEmpty(t, auditLogs[0].AdminName)
}

func TestGetAuditLogsUiFriendlyByAdmin(t *testing.T) {
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	admin := model.NewAdminUser("bob", uuid.New().String(), "uwu", "uwu")
	require.NoError(t, model.InsertAdminUser(tx, *admin))

	name := uuid.New().String()
	desc := uuid.New().String()

	auditLog := model.NewAuditLog(admin.Id, name, desc)
	require.NoError(t, model.InsertAuditLog(tx, auditLog))
	require.NoError(t, tx.Commit())

	auditLogs, err := model.GetAuditLogsUiFriendlyForAdmin(db, admin.Id)
	require.NoError(t, err)
	require.Len(t, auditLogs, 1)

	require.NotEmpty(t, auditLogs[0])
	require.NotEmpty(t, auditLogs[0].AdminName)
}

func TestGetAuditLogsUiFriendlyByGame(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	admin := model.NewAdminUser("bob", uuid.New().String(), "uwu", "uwu")
	require.NoError(t, model.InsertAdminUser(tx, *admin))

	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayerTx(tx, p1))

	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayerTx(tx, p2))

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))

	_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey, model.Score_Win, r)
	require.NoError(t, err)

	auditLog := model.NewAuditLog(admin.Id, "test-1", "test-2")
	require.NoError(t, model.InsertAuditLog(tx, auditLog))
	require.NoError(t, model.InsertAuditLogGameAffected(tx, &model.AuditLogGameAffected{
		AuditLogId: auditLog.Id,
		GameIkey:   ikey,
	}))

	require.NoError(t, tx.Commit())

	logs, err := model.GetAuditLogsUiFriendlyForGame(db, ikey)
	require.NoError(t, err)
	require.Len(t, logs, 1)

	require.Equal(t, auditLog.Id, logs[0].Id)
}

func TestGetAuditLogWithGameAndPlayer(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	admin := model.NewAdminUser("bob", uuid.New().String(), "uwu", "uwu")
	require.NoError(t, model.InsertAdminUser(tx, *admin))

	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayerTx(tx, p1))

	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayerTx(tx, p2))

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))

	_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey, model.Score_Win, r)
	require.NoError(t, err)

	auditLog := model.NewAuditLog(admin.Id, "test-1", "test-2")
	require.NoError(t, model.InsertAuditLog(tx, auditLog))

	gameAuditLog := model.AuditLogGameAffected{
		AuditLogId: auditLog.Id,
		GameIkey:   ikey,
	}
	require.NoError(t, model.InsertAuditLogGameAffected(tx, &gameAuditLog))

	p1AuditLog := model.AuditLogPlayerAffected{
		AuditLogId: auditLog.Id,
		PlayerId:   p1.Id,
		EloChange:  123,
	}
	require.NoError(t, model.InsertAuditLogPlayerAffected(tx, &p1AuditLog))

	p2AuditLog := model.AuditLogPlayerAffected{
		AuditLogId: auditLog.Id,
		PlayerId:   p2.Id,
		EloChange:  456,
	}
	require.NoError(t, model.InsertAuditLogPlayerAffected(tx, &p2AuditLog))

	log, err := model.GetAuditLog(tx, auditLog.Id)
	require.NoError(t, err)

	require.NoError(t, tx.Commit())

	require.Equal(t, auditLog.Id, log.Id)
	require.Len(t, log.Players, 2)
	require.Len(t, log.Games, 1)

	require.Contains(t, log.Games, gameAuditLog)
}
