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
	PlayerWhite        uuid.UUID `db:"player_white"`
	PlayerBlack        uuid.UUID `db:"player_black"`
	Score              Score     `db:"score"`
	Submitter          uuid.UUID `db:"submitter"`
	Played             time.Time `db:"played"`
	Deleted            bool      `db:"deleted"`
	DEPRECATEDEloGiven int       `db:"elo_given"` // Deprecated: for use with old elo system
	DEPRECATEDEloTaken int       `db:"elo_taken"` // Deprecated: for use with old elo system
	// Liglicko2White and Liglicko2Black are per-game liglicko2 deltas for each side.
	// They preserve sign, so draws between uneven players can still show non-zero
	// changes.
	Liglicko2White  float64 `db:"liglicko2_white"`
	Liglicko2Black  float64 `db:"liglicko2_black"`
	SubmitIp        string  `db:"submit_ip"`
	SubmitUserAgent string  `db:"submit_user_agent"`
	IKey            int64   `db:"ikey"`
}

type GameWithPlayerNames struct {
	Game
	WhiteName string `db:"white_name"`
	BlackName string `db:"black_name"`
}

type GameWithOutcome struct {
	Ikey            int64
	PlayerName      string
	OpponentName    string
	Outcome         string
	Played          time.Time
	EloChange       int
	Liglicko2Change float64
}

type GameWithDetails struct {
	Game
	WhiteName     string `db:"white_name"`
	BlackName     string `db:"black_name"`
	SubmitterName string `db:"submitter_name"`
}

func (g GameWithDetails) WinnerName() string {
	if g.Score == Score_Win {
		return g.WhiteName
	} else if g.Score == Score_Loss {
		return g.BlackName
	}
	return "Draw"
}

func GetGamesWithOutcomes(db *db.Db) ([]GameWithOutcome, error) {
	var games []GameWithPlayerNames
	err := db.GetSqlxDb().Select(&games, `
SELECT g.*, w.name as white_name, b.name as black_name
FROM games g
JOIN players w ON g.player_white = w.id
JOIN players b ON g.player_black = b.id
ORDER BY g.played DESC;`)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get games by player"), err)
	}

	gamesWithOutcomes := make([]GameWithOutcome, 0)
	for _, g := range games {
		gamesWithOutcomes = append(gamesWithOutcomes, g.MapGameToGameWithOutcome(g.PlayerWhite))
	}

	return gamesWithOutcomes, nil
}

func GetGameWithDetails(db *db.Db, ikey int64) (GameWithDetails, error) {
	var game GameWithDetails
	err := db.GetSqlxDb().Get(&game, `
SELECT g.*, w.name as white_name, b.name as black_name, s.name as submitter_name
FROM games g
JOIN players w ON g.player_white = w.id
JOIN players b ON g.player_black = b.id
JOIN players s ON g.submitter = s.id
WHERE g.ikey = $1;`, ikey)
	if err != nil {
		return GameWithDetails{}, errors.Join(errors.New("Cannot get game details"), err)
	}

	return game, nil
}

type GamesUiFriendly struct {
	Wins, Draws, Losses         int
	TotalGames                  int
	WinRate, LossRate, DrawRate float64
	WhiteWinRate, BlackWinRate  float64
	Games                       []GameWithOutcome
}

func (g *GameWithPlayerNames) MapGameToGameWithOutcome(playerId uuid.UUID) GameWithOutcome {
	gw := GameWithOutcome{
		Played: g.Played,
		Ikey:   g.IKey,
	}

	isWhite := g.PlayerWhite == playerId
	if isWhite {
		gw.OpponentName = g.BlackName
		gw.PlayerName = g.WhiteName
		if g.Score == Score_Win {
			gw.Outcome = "Win"
			gw.EloChange = g.DEPRECATEDEloGiven
			gw.Liglicko2Change = g.Liglicko2White
		} else if g.Score == Score_Loss {
			gw.Outcome = "Loss"
			gw.EloChange = g.DEPRECATEDEloTaken
			gw.Liglicko2Change = g.Liglicko2White
		} else {
			gw.Outcome = "Draw"
			gw.EloChange = 0
			gw.Liglicko2Change = g.Liglicko2White
		}
	} else {
		gw.OpponentName = g.WhiteName
		gw.PlayerName = g.BlackName
		if g.Score == Score_Loss {
			gw.Outcome = "Win"
			gw.EloChange = g.DEPRECATEDEloGiven
			gw.Liglicko2Change = g.Liglicko2Black
		} else if g.Score == Score_Win {
			gw.Outcome = "Loss"
			gw.EloChange = g.DEPRECATEDEloTaken
			gw.Liglicko2Change = g.Liglicko2Black
		} else {
			gw.Outcome = "Draw"
			gw.EloChange = 0
			gw.Liglicko2Change = g.Liglicko2Black
		}
	}

	return gw
}

func MapGamesToUserFriendly(playerId uuid.UUID, games []GameWithPlayerNames) GamesUiFriendly {
	details := GamesUiFriendly{
		Games:      make([]GameWithOutcome, 0),
		TotalGames: len(games),
	}

	var whiteGames, whiteWins, blackGames, blackWins int
	for _, g := range games {

		isWhite := g.PlayerWhite == playerId
		if isWhite {
			whiteGames++
			if g.Score == Score_Win {
				details.Wins++
				whiteWins++
			} else if g.Score == Score_Loss {
				details.Losses++
			} else {
				details.Draws++
			}
		} else {
			blackGames++
			if g.Score == Score_Loss {
				details.Wins++
				blackWins++
			} else if g.Score == Score_Win {
				details.Losses++
			} else {
				details.Draws++
			}
		}

		gw := g.MapGameToGameWithOutcome(playerId)
		details.Games = append(details.Games, gw)
	}

	if details.TotalGames > 0 {
		details.WinRate = float64(details.Wins) / float64(details.TotalGames) * 100
		details.LossRate = float64(details.Losses) / float64(details.TotalGames) * 100
		details.DrawRate = float64(details.Draws) / float64(details.TotalGames) * 100
		details.WhiteWinRate = float64(whiteWins) / float64(whiteGames) * 100
		details.BlackWinRate = float64(blackWins) / float64(blackGames) * 100
	}

	return details
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

func CreateGame(tx *sqlx.Tx, submitter, opponent *Player, submitterIsWhite bool, ikey int64, score Score, r *http.Request) (Game, int, int, error) {
	if submitter.Id == opponent.Id {
		return Game{}, 0, 0, errors.New("Both players are the same")
	}

	game := Game{
		Score:           score,
		Submitter:       submitter.Id,
		Played:          time.Now(),
		Deleted:         false,
		SubmitIp:        r.RemoteAddr,
		SubmitUserAgent: r.UserAgent(),
		IKey:            ikey,
	}

	var pWhite, pBlack *Player
	if submitterIsWhite {
		pWhite = submitter
		pBlack = opponent
	} else {
		pWhite = opponent
		pBlack = submitter
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

	eloWhite, eloBlack := CalculateElo(pWhite, pBlack, outcome)
	liglicko2White, liglicko2Black, err := CalculateLiglicko2(pWhite, pBlack, outcome, game.Played)
	if err != nil {
		return Game{}, 0, 0, errors.Join(errors.New("Could not calculate liglicko2"), err)
	}
	game.PlayerWhite = pWhite.Id
	game.PlayerBlack = pBlack.Id

	if eloWhite > eloBlack {
		game.DEPRECATEDEloGiven = eloWhite
		game.DEPRECATEDEloTaken = eloBlack
	} else {
		game.DEPRECATEDEloGiven = eloBlack
		game.DEPRECATEDEloTaken = eloWhite
	}

	game.Liglicko2White = liglicko2White
	game.Liglicko2Black = liglicko2Black

	_, err = tx.NamedExec(`
INSERT INTO games (player_white, player_black, score, submitter, played, deleted, elo_given, elo_taken, liglicko2_white, liglicko2_black, submit_ip, submit_user_agent, ikey)
VALUES (:player_white, :player_black, :score, :submitter, :played, :deleted, :elo_given, :elo_taken, :liglicko2_white, :liglicko2_black, :submit_ip, :submit_user_agent, :ikey);
  	`, game)

	if err != nil {
		return Game{}, 0, 0, errors.Join(errors.New("Cannot insert game"), err)
	}

	_, err = tx.NamedExec(`UPDATE players 
SET elo=:elo, liglicko2_rating=:liglicko2_rating, liglicko2_deviation=:liglicko2_deviation, liglicko2_volatility=:liglicko2_volatility, liglicko2_at=:liglicko2_at
WHERE id=:id`, pWhite)
	if err != nil {
		return Game{}, 0, 0, errors.Join(errors.New("Cannot set elo of white player"), err)
	}

	_, err = tx.NamedExec(`UPDATE players 
SET elo=:elo, liglicko2_rating=:liglicko2_rating, liglicko2_deviation=:liglicko2_deviation, liglicko2_volatility=:liglicko2_volatility, liglicko2_at=:liglicko2_at
WHERE id=:id`, pBlack)
	if err != nil {
		return Game{}, 0, 0, errors.Join(errors.New("Cannot set elo of black player"), err)
	}

	return game, int(liglicko2White), int(liglicko2Black), nil
}

func SubmitGame(db *db.Db, whiteName, blackName string, submitterIsWhite bool, ikey int64, score Score, r *http.Request) (*Game, *Player, *Player, int, int, error) {
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	if err != nil {
		return nil, nil, nil, 0, 0, errors.Join(errors.New("Could not start transaction"), err)
	}
	defer tx.Rollback()

	white, err := getOrCreatePlayer(tx, whiteName)
	if err != nil {
		return nil, nil, nil, 0, 0, errors.Join(errors.New("Could not get or create white player"), err)
	}

	black, err := getOrCreatePlayer(tx, blackName)
	if err != nil {
		return nil, nil, nil, 0, 0, errors.Join(errors.New("Could not get or create black player"), err)
	}

	var submitter, opponent *Player
	if submitterIsWhite {
		submitter = &white
		opponent = &black
	} else {
		submitter = &black
		opponent = &white
	}

	game, eloWhite, eloBlack, err := CreateGame(tx, submitter, opponent, submitterIsWhite, ikey, score, r)
	if err != nil {
		return nil, nil, nil, 0, 0, errors.Join(errors.New("Could not create game"), err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, nil, nil, 0, 0, errors.Join(errors.New("Could not commit transaction"), err)
	}

	return &game, &white, &black, eloWhite, eloBlack, nil
}

func GetGamesByPlayer(db *db.Db, playerId uuid.UUID) ([]GameWithPlayerNames, error) {
	var games []GameWithPlayerNames
	err := db.GetSqlxDb().Select(&games, `
SELECT g.*, w.name as white_name, b.name as black_name
FROM games g
JOIN players w ON g.player_white = w.id
JOIN players b ON g.player_black = b.id
WHERE (g.player_white=$1 OR g.player_black=$1) AND g.deleted=false
ORDER BY g.played DESC`, playerId)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get games by player"), err)
	}

	return games, nil
}

func GetTotalGameCount(db *db.Db) (int, error) {
	var count int
	err := db.GetSqlxDb().Get(&count, "SELECT count(*) FROM games WHERE deleted=false")
	if err != nil {
		return 0, errors.Join(errors.New("Cannot get total game count"), err)
	}

	return count, nil
}
