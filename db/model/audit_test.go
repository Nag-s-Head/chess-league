package model_test

import (
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
