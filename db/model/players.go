package model

import (
	"time"

	"github.com/google/uuid"
)

type Player struct {
	Id             uuid.UUID `db:"id"`
	Name           string    `db:"name"`
	NameNormalised string    `db:"name_normalised"`
	Elo            int       `db:"elo"`
	JoinTime       time.Time `db:"join_time"`
}
