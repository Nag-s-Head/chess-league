package submitgame

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/rules"
	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	BasePath          = "/submit-game"
	IKeyCookie        = "ikey"
	MagicNumberParam  = "magic"
	MagicNumberCookie = "magic"
	MagicNumberEnvVar = "MAGIC_NUMBER"
)

type PlayerConsolidationModel struct {
	Results      []PlayerLookupResult
	PlayersGames []model.Game
}

type PlayerLookupResult struct {
	Players        []model.Player
	ExactMatch     bool
	IsWhite        bool
	Name           string
	NameNormalised string
}

var magicNumber string = os.Getenv(MagicNumberEnvVar)

func VerifyMagic(r *http.Request) bool {
	cookie, err := r.Cookie(MagicNumberCookie)
	if err == nil && cookie.Value != "" {
		return cookie.Value == magicNumber
	}
	return r.URL.Query().Get(MagicNumberParam) == magicNumber
}

func GetLookupResult(db db.Db, name string, isWhite bool) (PlayerLookupResult, error) {
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

func GetGamesForPairs(db db.Db, playersA PlayerLookupResult, playersB PlayerLookupResult) ([]model.Game, error) {
	playerAIds := make([]uuid.UUID, 0)
	playerBIds := make([]uuid.UUID, 0)

	for _, player := range playersA.Players {
		playerAIds = append(playerAIds, player.Id)
	}

	for _, player := range playersB.Players {
		playerBIds = append(playerBIds, player.Id)
	}

	allGameCombs, err := model.GetGamesByPlayerPairCombs(db, playerAIds, playerBIds)

	if err != nil {
		return make([]model.Game, 0), errors.Join(errors.New("Cannot find games for pair"), err)
	}

	return allGameCombs, nil
}

const (
	playerName      = "player-name"
	playedAs        = "played-as"
	otherPlayerName = "other-player-name"
	winner          = "winner"

	whitePlayerName = "white-player-name"
	blackPlayerName = "black-player-name"
)

func renderUserLookup(db db.Db, player1, player2 string, player1White bool) ([]byte, error) {
	// Lookup the players
	results := new(PlayerConsolidationModel)

	res, err := GetLookupResult(db, player1, player1White)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot lookup player 1"), err)
	}
	results.Results = append(results.Results, res)

	res, err = GetLookupResult(db, player2, !player1White)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot lookup player 2"), err)
	}
	results.Results = append(results.Results, res)

	allExact := true
	for _, res := range results.Results {
		allExact = allExact && res.ExactMatch
	}

	games, err := GetGamesForPairs(db, results.Results[0], results.Results[1])
	if err != nil {
		return nil, errors.Join(errors.New("Cannot lookup games"), err)
	}
	results.PlayersGames = append(results.PlayersGames, games...)

	// Check results and send to the UI
	var buf bytes.Buffer

	err = resultsTpl.Execute(&buf, results)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot execute template"), err)
	}

	return buf.Bytes(), nil
}

func doUserLookupSubmit(db db.Db, w http.ResponseWriter, r *http.Request) error {
	player1 := strings.TrimSpace(r.FormValue(playerName))
	if player1 == "" {
		return errors.New("Player 1 is not set")
	}

	player2 := strings.TrimSpace(r.FormValue(otherPlayerName))
	if player2 == "" {
		return errors.New("Player 2 is not set")
	}

	rawPlayedAs := r.FormValue(playedAs)
	if !(rawPlayedAs == "white" || rawPlayedAs == "black") {
		return fmt.Errorf("Played as value %s is not valid", rawPlayedAs)
	}
	player1White := rawPlayedAs == "white"

	buf, err := renderUserLookup(db, player1, player2, player1White)
	if err != nil {
		return err
	}

	w.Write(buf)
	return nil
}

func doFinalSubmit(db db.Db, w http.ResponseWriter, r *http.Request) error {
	white := strings.TrimSpace(r.FormValue(whitePlayerName))
	black := strings.TrimSpace(r.FormValue(blackPlayerName))

	rawPlayedAs := r.FormValue(playedAs)
	submitterIsWhite := rawPlayedAs == "white"

	ikeyCookie, err := r.Cookie(IKeyCookie)
	if err != nil {
		return errors.Join(errors.New("Could not find ikey cookie"), err)
	}

	ikey, err := strconv.ParseInt(ikeyCookie.Value, 10, 64)
	if err != nil {
		return errors.New("Could not read ikey cookie, please refresh the page")
	}

	var score model.Score
	winner := r.FormValue(winner)
	if winner == "win" {
		if submitterIsWhite {
			score = model.Score_Win
		} else {
			score = model.Score_Loss
		}
	} else if winner == "loss" {
		if submitterIsWhite {
			score = model.Score_Loss
		} else {
			score = model.Score_Win
		}
	} else if winner == "draw" {
		score = model.Score_Draw
	} else {
		return errors.New("Invalid winner")
	}

	game, playerWhite, playerBlack, eloWhite, eloBlack, err := model.SubmitGame(db, white, black, submitterIsWhite, ikey, score, r)
	if err != nil {
		return errors.Join(errors.New("Could not submit game"), err)
	}

	slog.Info("Submitted a game", "game", game, "playerWhite", playerWhite, "playerBlack", playerBlack, "ip", model.GetRemoteAddr(r))

	http.SetCookie(w, &http.Cookie{
		Name:   IKeyCookie,
		Value:  "",
		MaxAge: 0,
		Path:   BasePath,
	})

	type FinalPlayer struct {
		EloGiven int
		model.Player
	}

	var buf bytes.Buffer
	err = successTpl.Execute(&buf, []FinalPlayer{
		{
			EloGiven: eloWhite,
			Player:   *playerWhite,
		},
		{
			EloGiven: eloBlack,
			Player:   *playerBlack,
		},
	})
	if err != nil {
		return err
	}

	w.Write(buf.Bytes())
	return nil
}

func DoSubmit(db db.Db, w http.ResponseWriter, r *http.Request) error {
	if !rules.HasAgreedToRules(r) {
		return errors.New("You must agree to the rules before submitting a game")
	}

	if !VerifyMagic(r) {
		slog.Warn("An attempt to access submit the form without the magic number was made", "ip", model.GetRemoteAddr(r))
		return errors.New("Magic number for submit is invalid")
	}

	err := r.ParseForm()
	if err != nil {
		return errors.Join(errors.New("Could not parse form"), err)
	}

	submitType := r.FormValue("submit-type")

	if submitType == "final" {
		slog.Info("Submitting a game", "form", r.Form, "ip", model.GetRemoteAddr(r))
		return doFinalSubmit(db, w, r)
	} else {
		slog.Info("Doing user lookup", "form", r.Form, "ip", model.GetRemoteAddr(r))
		return doUserLookupSubmit(db, w, r)
	}
}

type Error struct {
	Error string
}

type WsMessage struct {
	Player1Name string `json:"player-name"`
	Player2Name string `json:"other-player-name"`
	PlayedAs    string `json:"played-as"`
	SubmitType  string `json:"submit-type"`
}

func HandleWs(db db.Db, w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	if !rules.HasAgreedToRules(r) {
		return
	}

	if !VerifyMagic(r) {
		slog.Warn("An attempt to access submit the form without the magic number was made", "ip", model.GetRemoteAddr(r))
		return
	}

	const bufSize = 1024
	upgrader := websocket.Upgrader{
		HandshakeTimeout:  time.Second,
		ReadBufferSize:    bufSize,
		WriteBufferSize:   bufSize,
		EnableCompression: true,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Could start lookup game websocket", "err", err, "ip", model.GetRemoteAddr(r))

		var buf bytes.Buffer
		err = errorTpl.Execute(&buf, Error{Error: err.Error()})
		if err != nil {
			return
		}
	}

	defer conn.Close()

	var lastMessage WsMessage
	for {
		messageType, messageBytes, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Could not recieve from ws", "err", err)
			return
		}

		if messageType == websocket.PingMessage {
			err := conn.WriteMessage(websocket.PongMessage, []byte{})
			if err != nil {
				slog.Error("Could not send pong ws message", "err", err)
				return
			}

			continue
		}

		slog.Debug("Recieved a message", "message", string(messageBytes))

		var message WsMessage
		err = json.Unmarshal(messageBytes, &message)
		if err != nil {
			slog.Error("Could not read ws message", "err", err)
			return
		}

		if message.SubmitType == "final" {
			continue
		}

		// Prevents sending excess data
		if message == lastMessage {
			continue
		}
		lastMessage = message

		player1 := strings.TrimSpace(message.Player1Name)
		if player1 == "" {
			continue
		}

		player2 := strings.TrimSpace(message.Player2Name)
		if player2 == "" {
			continue
		}

		player1White := message.PlayedAs == "white"
		res, err := renderUserLookup(db, player1, player2, player1White)
		if err != nil {
			slog.Debug("There was an error whilst processing partial form data", "err", err)
		}

		if res != nil {
			payload := fmt.Sprintf(`<div class="flex flex-col gap-5 w-full" id="response">%s</div>`, string(res))
			err := conn.WriteMessage(websocket.TextMessage, []byte(payload))
			if err != nil {
				slog.Error("Could not send ws message", "err", err)
				return
			}
		}
	}
}

func Register(mux *http.ServeMux, db db.Db) {
	mux.HandleFunc(fmt.Sprintf("GET %s/submit/ws", BasePath), func(w http.ResponseWriter, r *http.Request) {
		HandleWs(db, w, r)
	})

	mux.HandleFunc(fmt.Sprintf("POST %s/submit", BasePath), func(w http.ResponseWriter, r *http.Request) {
		err := DoSubmit(db, w, r)
		if err != nil {
			slog.Error("Could not submit a game", "err", err, "params", r.Form, "ip", model.GetRemoteAddr(r))

			var buf bytes.Buffer
			err := errorTpl.Execute(&buf, Error{Error: err.Error()})
			if err != nil {
				return
			}
			w.Write(buf.Bytes())
		}
	})
}
