package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	Id                    uuid.UUID `db:"id"`
	Created               time.Time `db:"created"`
	DoneBy                uuid.UUID `db:"done_by"`
	OperationName         string    `db:"operation_name"`
	OperationDescrription string    `db:"operation_description"`
}

type AuditLogPlayerAffected struct {
	AuditLogId uuid.UUID `db:"audit_log_id"`
	PlayerId   uuid.UUID `db:"player_id"`
	EloChange  int       `db:"elo_change"`
}
