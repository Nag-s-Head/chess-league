package model_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestScore_Switch(t *testing.T) {
	t.Parallel()

	win := model.Score_Win
	win.Switch()
	require.Equal(t, model.Score_Loss, win)

	loss := model.Score_Loss
	loss.Switch()
	require.Equal(t, model.Score_Win, loss)

	draw := model.Score_Draw
	draw.Switch()
	require.Equal(t, model.Score_Draw, draw)
}

func TestSwapGameWinner(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	// 1. Setup Admin and Players
	admin := model.NewAdminUser("Admin-Swap", uuid.New().String(), "127.0.0.1", "Go-Test")
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

	// 2. Create a game (P1 wins as White)
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	ikey, _ := model.NextIKey(db)
	_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey, model.Score_Win, r)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	// Capture ratings after game
	p1AfterGame, _ := model.GetPlayer(db, p1.Id)
	p2AfterGame, _ := model.GetPlayer(db, p2.Id)

	// 3. Swap winner (P2 wins as Black now)
	err = model.SwapGameWinner(db, adminId, ikey)
	require.NoError(t, err)

	// 4. Verify Swap
	gameAfter, err := model.GetGameWithDetails(db, ikey)
	require.NoError(t, err)
	require.Equal(t, model.Score_Loss, gameAfter.Score) // P1 (White) is now Score_Loss

	// 5. Verify Ratings Replayed
	p1AfterSwap, _ := model.GetPlayer(db, p1.Id)
	p2AfterSwap, _ := model.GetPlayer(db, p2.Id)

	require.NotEqual(t, p1AfterGame.Liglicko2Rating, p1AfterSwap.Liglicko2Rating)
	require.NotEqual(t, p2AfterGame.Liglicko2Rating, p2AfterSwap.Liglicko2Rating)

	// Since it was 1-0 and now 0-1, P1's rating should be lower than before swap
	require.Less(t, p1AfterSwap.Liglicko2Rating, p1AfterGame.Liglicko2Rating)

	// 6. Verify Audit Logs
	var auditLogs []model.AuditLog
	err = db.GetSqlxDb().Select(&auditLogs, "SELECT * FROM audit_logs WHERE done_by = $1 ORDER BY created DESC", adminId)
	require.NoError(t, err)
	require.NotEmpty(t, auditLogs)
	require.Equal(t, "Swap Game Winner", auditLogs[0].OperationName)
	require.Contains(t, auditLogs[0].OperationDescription, fmt.Sprintf("Swapped winner for game %d", ikey))
}

func TestSetGameToDraw(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	// 1. Setup Admin and Players
	admin := model.NewAdminUser("Admin-Draw", uuid.New().String(), "127.0.0.1", "Go-Test")
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

	// 2. Create a game (P1 wins as White)
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	ikey, _ := model.NextIKey(db)
	_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey, model.Score_Win, r)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	// Capture ratings after game
	p1AfterGame, _ := model.GetPlayer(db, p1.Id)

	// 3. Set to Draw
	err = model.SetGameToDraw(db, adminId, ikey)
	require.NoError(t, err)

	// 4. Verify Change
	gameAfter, err := model.GetGameWithDetails(db, ikey)
	require.NoError(t, err)
	require.Equal(t, model.Score_Draw, gameAfter.Score)

	// 5. Verify Ratings Replayed
	p1AfterDraw, _ := model.GetPlayer(db, p1.Id)
	require.NotEqual(t, p1AfterGame.Liglicko2Rating, p1AfterDraw.Liglicko2Rating)

	// 6. Verify Audit Logs
	var auditLogs []model.AuditLog
	err = db.GetSqlxDb().Select(&auditLogs, "SELECT * FROM audit_logs WHERE done_by = $1 ORDER BY created DESC", adminId)
	require.NoError(t, err)
	require.NotEmpty(t, auditLogs)
	require.Equal(t, "Game Set To Draw", auditLogs[0].OperationName)
	require.Contains(t, auditLogs[0].OperationDescription, fmt.Sprintf("Set game %d to be a draw", ikey))
}

func TestSwapGameWinner_NonExistent(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	adminId := uuid.New()
	err := model.SwapGameWinner(db, adminId, 999999)
	require.Error(t, err)
}

func TestSetGameToDraw_NonExistent(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	adminId := uuid.New()
	err := model.SetGameToDraw(db, adminId, 999999)
	require.Error(t, err)
}
