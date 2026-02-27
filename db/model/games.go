package model

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Score string

const (
	Score_Win  Score = "1-0"
	Score_Loss Score = "0-1"
	Score_Draw Score = "1/2-1/2"
)

type Game struct {
	PlayerWhite     uuid.UUID `db:"player_white"`
	PlayerBlack     uuid.UUID `db:"player_black"`
	Score           Score     `db:"score"`
	Submitter       uuid.UUID `db:"submitter"`
	Played          time.Time `db:"played"`
	Deleted         bool      `db:"deleted"`
	EloGiven        int       `db:"elo_given"`
	EloTaken        int       `db:"elo_taken"`
	SubmitIp        string    `db:"submit_ip"`
	SubmitUserAgent string    `db:"submit_user_agent"`
	IKey            int64     `db:"ikey"`
}

func NextIKey(db *db.Db) (int64, error) {
	var ikey int64
	row := db.GetSqlxDb().QueryRow("SELECT nextval('game_ikey_sequence');")
	err := row.Scan(&ikey)
	if err != nil {
		return 0, errors.Join(errors.New("Cannot create new ikey"), err)
	}

	return ikey, nil
}

func CreateGame(tx *sqlx.Tx, player1, player2 *Player, isWhite bool, ikey int64, score Score, r *http.Request) (Game, error) {
	if player1.Id == player2.Id {
		return Game{}, errors.New("Both players are the same")
	}

	game := Game{
		Score:           score,
		Submitter:       player1.Id,
		Played:          time.Now(),
		Deleted:         false,
		SubmitIp:        r.RemoteAddr,
		SubmitUserAgent: r.UserAgent(),
		IKey:            ikey,
	}

	if !isWhite {
		tmp := player1
		player1 = player2
		player2 = tmp
	}

	var outcome Outcome
	switch score {
	case Score_Win:
		outcome = Outcome_Win
	case Score_Loss:
		outcome = Outcome_Loss
	case Score_Draw:
		outcome = Outcome_Draw
	}

	eloA, eloB := CalculateElo(player1, player2, outcome)
	game.PlayerWhite = player1.Id
	game.PlayerBlack = player2.Id

	if eloA > 0 {
		game.EloGiven = eloA
		game.EloTaken = eloB
	} else {
		game.EloGiven = eloB
		game.EloTaken = eloA
	}

	_, err := tx.NamedExec(`
INSERT INTO games (player_white, player_black, score, submitter, played, deleted, elo_given, elo_taken, submit_ip, submit_user_agent, ikey)
VALUES (:player_white, :player_black, :score, :submitter, :played, :deleted, :elo_given, :elo_taken, :submit_ip, :submit_user_agent, :ikey);
  	`, game)

	if err != nil {
		return Game{}, errors.Join(errors.New("Cannot insert game"), err)
	}

	_, err = tx.NamedExec(`UPDATE players SET elo=:elo WHERE id=:id`, player1)
	if err != nil {
		return Game{}, errors.Join(errors.New("Cannot set elo of player 1"), err)
	}

	_, err = tx.NamedExec(`UPDATE players SET elo=:elo WHERE id=:id`, player2)
	if err != nil {
		return Game{}, errors.Join(errors.New("Cannot set elo of player 2"), err)
	}

	return game, nil
}

func SubmitGame(db *db.Db, p1Name, p2Name string, isWhite bool, ikey int64, score Score, r *http.Request) (*Game, *Player, *Player, error) {
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	if err != nil {
		return nil, nil, nil, errors.Join(errors.New("Could not start transaction"), err)
	}
	defer tx.Rollback()

	player1, err := getOrCreatePlayer(tx, p1Name)
	if err != nil {
		return nil, nil, nil, errors.Join(errors.New("Could not get or create player 1"), err)
	}

	player2, err := getOrCreatePlayer(tx, p2Name)
	if err != nil {
		return nil, nil, nil, errors.Join(errors.New("Could not get or create player 2"), err)
	}

	game, err := CreateGame(tx, &player1, &player2, isWhite, ikey, score, r)
	if err != nil {
		return nil, nil, nil, errors.Join(errors.New("Could not create game"), err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, nil, nil, errors.Join(errors.New("Could not commit transaction"), err)
	}

	return &game, &player1, &player2, nil
}
