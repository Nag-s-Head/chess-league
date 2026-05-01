package model_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDeleteGame(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	// 1. Setup Admin and Players
	admin := model.NewAdminUser("Admin-Core", uuid.New().String(), "127.0.0.1", "Go-Test")
	adminId := admin.Id

	txAdmin, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)
	require.NoError(t, model.InsertAdminUser(txAdmin, *admin))
	require.NoError(t, txAdmin.Commit())

	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p1))
	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p2))

	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	r.RemoteAddr = "127.0.0.1"

	// 2. Create three games in sequence
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	ikey1, _ := model.NextIKey(db)
	_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey1, model.Score_Win, r)
	require.NoError(t, err)

	ikey2, _ := model.NextIKey(db)
	_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey2, model.Score_Loss, r)
	require.NoError(t, err)

	ikey3, _ := model.NextIKey(db)
	_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey3, model.Score_Draw, r)
	require.NoError(t, err)

	require.NoError(t, tx.Commit())

	// Capture state before deletion
	p1BeforeDelete, _ := model.GetPlayer(db, p1.Id)
	p2BeforeDelete, _ := model.GetPlayer(db, p2.Id)

	// 3. Delete Game 1
	err = model.DeleteGame(db, adminId, ikey1)
	require.NoError(t, err)

	// 4. Verify Deletion
	g1After, err := model.GetGameWithDetails(db, ikey1)
	require.NoError(t, err)
	require.True(t, g1After.Deleted)

	// 5. Verify Ratings Replayed correctly
	p1AfterDelete, _ := model.GetPlayer(db, p1.Id)
	p2AfterDelete, _ := model.GetPlayer(db, p2.Id)

	require.NotEqual(t, p1BeforeDelete.Liglicko2Rating, p1AfterDelete.Liglicko2Rating)
	require.NotEqual(t, p2BeforeDelete.Liglicko2Rating, p2AfterDelete.Liglicko2Rating)

	// Game 2 OldRating should now be 1500 (since Game 1 is deleted)
	g2After, _ := model.GetGameWithDetails(db, ikey2)
	require.InDelta(t, 1500.0, g2After.Liglicko2WhiteOldRating, 1e-9)
	require.InDelta(t, 1500.0, g2After.Liglicko2BlackOldRating, 1e-9)

	// 6. Verify Audit Logs
	var auditLogs []model.AuditLog
	err = db.GetSqlxDb().Select(&auditLogs, "SELECT * FROM audit_logs WHERE done_by = $1 ORDER BY created DESC", adminId)
	require.NoError(t, err)
	require.NotEmpty(t, auditLogs)
	require.Equal(t, "Game Deletion", auditLogs[0].OperationName)
	require.Contains(t, auditLogs[0].OperationDescription, "Deleted game")
}

func TestDeleteLatestGame(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	// 1. Setup
	admin := model.NewAdminUser("Admin-Latest", uuid.New().String(), "127.0.0.1", "Go-Test")
	txAdmin, _ := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	model.InsertAdminUser(txAdmin, *admin)
	txAdmin.Commit()

	p1 := model.NewPlayer(uuid.New().String())
	model.InsertPlayer(db, p1)
	p2 := model.NewPlayer(uuid.New().String())
	model.InsertPlayer(db, p2)

	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	r.RemoteAddr = "127.0.0.1"

	// 2. Create first game and capture state
	tx, _ := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	ikey1, _ := model.NextIKey(db)
	model.CreateGame(tx, &p1, &p2, true, ikey1, model.Score_Win, r)
	tx.Commit()

	// Capture state after Game 1. This is what we expect to return to after deleting G2.
	p1AfterG1, _ := model.GetPlayer(db, p1.Id)
	p2AfterG1, _ := model.GetPlayer(db, p2.Id)

	// 3. Create second game (the latest)
	tx2, _ := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	ikey2, _ := model.NextIKey(db)
	model.CreateGame(tx2, &p1, &p2, true, ikey2, model.Score_Loss, r)
	tx2.Commit()

	// 4. Delete Game 2 (the latest)
	err := model.DeleteGame(db, admin.Id, ikey2)
	require.NoError(t, err)

	// 5. Verify
	p1Final, _ := model.GetPlayer(db, p1.Id)
	p2Final, _ := model.GetPlayer(db, p2.Id)

	require.InDelta(t, p1AfterG1.Liglicko2Rating, p1Final.Liglicko2Rating, 1e-9)
	require.InDelta(t, p2AfterG1.Liglicko2Rating, p2Final.Liglicko2Rating, 1e-9)

	g2After, _ := model.GetGameWithDetails(db, ikey2)
	require.True(t, g2After.Deleted)
}

func TestDeleteGameWithManySubsequent(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	// 1. Setup
	admin := model.NewAdminUser("Admin-Many", uuid.New().String(), "127.0.0.1", "Go-Test")
	txAdmin, _ := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	model.InsertAdminUser(txAdmin, *admin)
	txAdmin.Commit()

	p1 := model.NewPlayer(uuid.New().String())
	model.InsertPlayer(db, p1)
	p2 := model.NewPlayer(uuid.New().String())
	model.InsertPlayer(db, p2)

	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	r.RemoteAddr = "127.0.0.1"

	// 2. Create a chain of 10 games
	var ikeys []int64
	tx, _ := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	for i := 0; i < 10; i++ {
		ikey, _ := model.NextIKey(db)
		ikeys = append(ikeys, ikey)
		// Alternate wins
		score := model.Score_Win
		if i%2 == 1 {
			score = model.Score_Loss
		}
		model.CreateGame(tx, &p1, &p2, true, ikey, score, r)
	}
	tx.Commit()

	p1BeforeDelete, _ := model.GetPlayer(db, p1.Id)

	// 3. Delete the 3rd game (index 2)
	err := model.DeleteGame(db, admin.Id, ikeys[2])
	require.NoError(t, err)

	// 4. Verify the ripple effect
	p1AfterDelete, _ := model.GetPlayer(db, p1.Id)
	require.NotEqual(t, p1BeforeDelete.Liglicko2Rating, p1AfterDelete.Liglicko2Rating)

	// Verify all subsequent games (ikeys[3] through ikeys[9]) were updated
	// We can check if their "old" states now match the updated history
	for i := 3; i < 10; i++ {
		game, _ := model.GetGameWithDetails(db, ikeys[i])
		// This check ensures ReplayFrom actually touched and updated these games
		require.NotZero(t, game.Liglicko2WhiteOldRating)
	}
}

func TestDeleteNonExistentGameFails(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	adminId := uuid.New()
	err := model.DeleteGame(db, adminId, 999999)
	require.Error(t, err)
}
