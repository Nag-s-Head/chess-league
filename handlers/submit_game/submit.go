package submitgame

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/djpiper28/rpg-book/common/normalisation"
)

const (
	BasePath          = "/submit-game"
	IKeyCookie        = "ikey"
	MagicNumberParam  = "magic"
	MagicNumberCookie = "magic"
	MagicNumberEnvVar = "MAGIC_NUMBER"
)

type PlayerConsolidationModel struct {
	Results []PlayerLookupResult
}

type PlayerLookupResult struct {
	Players        []model.Player
	ExactMatch     bool
	IsWhite        bool
	Name           string
	NameNormalised string
}

func GetLookupResult(db *db.Db, name string, isWhite bool) (PlayerLookupResult, error) {
	nameNormalised := normalisation.Normalise(name)
	players, err := model.SearchPlayerByName(db, nameNormalised)
	if err != nil {
		return PlayerLookupResult{}, errors.Join(errors.New("Cannot find player 1 by name"), err)
	}

	if len(players) == 1 {
		if nameNormalised == players[0].NameNormalised {
			return PlayerLookupResult{
				Players:        players,
				ExactMatch:     true,
				IsWhite:        isWhite,
				Name:           players[0].Name,
				NameNormalised: nameNormalised,
			}, nil
		}
	}

	return PlayerLookupResult{
		Players:        players,
		IsWhite:        isWhite,
		Name:           name,
		NameNormalised: nameNormalised,
	}, nil
}

func DoSubmit(db *db.Db, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return errors.Join(errors.New("Could not parse form"), err)
	}

	const (
		playerName      = "player-name"
		playedAs        = "played-as"
		otherPlayerName = "other-player-name"
		winner          = "winner"
	)

	player1 := r.FormValue(playerName)
	if player1 == "" {
		return errors.New("Player 1 is not set")
	}

	player2 := r.FormValue(otherPlayerName)
	if player2 == "" {
		return errors.New("Player 2 is not set")
	}

	rawPlayedAs := r.FormValue(playedAs)
	if !(rawPlayedAs == "white" || rawPlayedAs == "black") {
		return fmt.Errorf("Played as value %s is not valid", rawPlayedAs)
	}
	player1White := rawPlayedAs == "white"

	// Lookup the players
	results := PlayerConsolidationModel{Results: make([]PlayerLookupResult, 0)}
	res, err := GetLookupResult(db, player1, player1White)
	if err != nil {
		return errors.Join(errors.New("Cannot lookup player 1"), err)
	}
	results.Results = append(results.Results, res)

	res, err = GetLookupResult(db, player2, !player1White)
	if err != nil {
		return errors.Join(errors.New("Cannot lookup player 2"), err)
	}
	results.Results = append(results.Results, res)

	allExact := true
	for _, res := range results.Results {
		allExact = allExact && res.ExactMatch
	}

	// Check results and send to the UI
	if allExact {
		w.Write([]byte("TODO: all players were exact"))
	} else {
		var buf bytes.Buffer
		err := resultsTpl.Execute(&buf, results)
		if err != nil {
			return errors.Join(errors.New("Cannot execute template"), err)
		}

		w.Write(buf.Bytes())
	}

	return nil
}

func Register(mux *http.ServeMux, db *db.Db) {
	mux.HandleFunc(fmt.Sprintf("POST %s/submit", BasePath), func(w http.ResponseWriter, r *http.Request) {
		err := DoSubmit(db, w, r)
		if err != nil {
			slog.Error("Could not submit a game", "err", err, "params", r.Form)
		}
	})
}
