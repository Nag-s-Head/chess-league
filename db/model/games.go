package model

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	PlayerWhite     uuid.UUID `db:"player_white"`
	PlayerBlack     uuid.UUID `db:"player_white"`
	Score           string    `db:"score"`
	Submitter       uuid.UUID `db:"submitter"`
	Played          time.Time `db:"played"`
	Deleted         bool      `db:"deleted"`
	EloGiven        int       `db:"elo_given"`
	EloTake         int       `db:"elo_taken"`
	SubmitIp        string    `db:"submimt_ip"`
	SubmitUserAgent string    `db:"submit_user_agent"`
}
