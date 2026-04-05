package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuditLog struct {
	Id                   uuid.UUID `db:"id"`
	Created              time.Time `db:"created"`
	DoneBy               uuid.UUID `db:"done_by"`
	OperationName        string    `db:"operation_name"`
	OperationDescription string    `db:"operation_description"`
}

func NewAuditLog(DoneBy uuid.UUID, OperationName, OperationDescription string) *AuditLog {
	return &AuditLog{
		Id:                   uuid.New(),
		Created:              time.Now(),
		DoneBy:               DoneBy,
		OperationName:        OperationName,
		OperationDescription: OperationDescription,
	}
}

func InsertAuditLog(tx *sqlx.Tx, auditLog *AuditLog) error {
	_, err := tx.NamedExec(`
	  INSERT INTO audit_logs (id, operation_name, operation_description, done_by, created)
		VALUES (:id, :operation_name, :operation_description, :done_by, :created);
		`, auditLog)
	return err
}

type AuditLogPlayerAffected struct {
	AuditLogId uuid.UUID `db:"audit_log_id"`
	PlayerId   uuid.UUID `db:"player_id"`
	EloChange  int       `db:"elo_change"`
}

func NewAuditLogPlayerAffected(auditId uuid.UUID, PlayerId uuid.UUID, EloChange int) *AuditLogPlayerAffected {
	return &AuditLogPlayerAffected{
		AuditLogId: auditId,
		PlayerId:   PlayerId,
		EloChange:  EloChange,
	}
}

func InsertAuditLogPlayerAffected(tx *sqlx.Tx, playerAffected *AuditLogPlayerAffected) error {
	_, err := tx.NamedExec(`
	  INSERT INTO audit_log_player_affected (audit_log_id, player_id, elo_change) 
		VALUES(:audit_log_id, :player_id, :elo_change);
	`, playerAffected)
	return err
}

type DetailedAuditLog struct {
	AuditLog
	Players []AuditLogPlayerAffected
}

func GetAuditLog(tx *sqlx.Tx, id uuid.UUID) (*DetailedAuditLog, error) {
	var auditLog AuditLog
	err := tx.Get(&auditLog, "SELECT * FROM audit_logs WHERE id=$1;", id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get audit log"), err)
	}

	players := make([]AuditLogPlayerAffected, 0)
	rows, err := tx.Queryx(`
		SELECT audit_log_player_affected.* 
		FROM audit_log_player_affected 
		INNER JOIN audit_logs 
		  ON audit_logs.id=audit_log_player_affected.audit_log_id
		WHERE audit_log_player_affected.audit_log_id=$1;`, id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get audit log player affected"), err)
	}

	defer rows.Close()

	for rows.Next() {
		var player AuditLogPlayerAffected
		err := rows.StructScan(&player)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot scan audit log player affacted"), err)
		}

		players = append(players, player)
	}

	result := &DetailedAuditLog{
		AuditLog: auditLog,
		Players:  players,
	}

	return result, nil
}
