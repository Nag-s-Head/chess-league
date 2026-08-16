package model

import (
	"errors"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
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
	AuditLogId   uuid.UUID `db:"audit_log_id"`
	PlayerId     uuid.UUID `db:"player_id"`
	IsMainTarget bool      `db:"is_main_target"`
	PlayerName   string    `db:"player_name"`
}

func NewAuditLogPlayerAffected(auditId uuid.UUID, PlayerId uuid.UUID, isMainTarget bool) *AuditLogPlayerAffected {
	return &AuditLogPlayerAffected{
		AuditLogId:   auditId,
		PlayerId:     PlayerId,
		IsMainTarget: isMainTarget,
	}
}

func InsertAuditLogPlayerAffected(tx *sqlx.Tx, playerAffected *AuditLogPlayerAffected) error {
	_, err := tx.NamedExec(`
	  INSERT INTO audit_log_player_affected (audit_log_id, player_id, is_main_target)
		VALUES(:audit_log_id, :player_id, :is_main_target);
	`, playerAffected)
	return err
}

type AuditLogGameAffected struct {
	AuditLogId   uuid.UUID `db:"audit_log_id"`
	GameIkey     int64     `db:"game_ikey"`
	IsMainTarget bool      `db:"is_main_target"`
	WhiteName    string    `db:"white_name"`
	BlackName    string    `db:"black_name"`
	Played       time.Time `db:"played"`
}

func InsertAuditLogGameAffected(tx *sqlx.Tx, gameAffected *AuditLogGameAffected) error {
	_, err := tx.NamedExec(`
	  INSERT INTO audit_log_game_affected (audit_log_id, game_ikey, is_main_target)
		VALUES(:audit_log_id, :game_ikey, :is_main_target);
	`, gameAffected)
	return err
}

type DetailedAuditLog struct {
	AuditLog
	AdminName string `db:"admin_name"`
	Players   []AuditLogPlayerAffected
	Games     []AuditLogGameAffected
}

func GetAuditLog(tx *sqlx.Tx, id uuid.UUID) (*DetailedAuditLog, error) {
	var auditLog struct {
		AuditLog
		AdminName string `db:"admin_name"`
	}
	err := tx.Get(&auditLog, `
		SELECT audit_logs.*, admin_users.name as admin_name 
		FROM audit_logs 
		JOIN admin_users ON audit_logs.done_by = admin_users.id
		WHERE audit_logs.id=$1;`, id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get audit log"), err)
	}

	players := make([]AuditLogPlayerAffected, 0)
	err = tx.Select(&players, `
		SELECT a.*, p.name as player_name 
		FROM audit_log_player_affected a
		JOIN players p ON a.player_id = p.id
		WHERE a.audit_log_id=$1
		ORDER BY is_main_target ASC, p.name ASC, p.id ASC;`, id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get audit log player affected"), err)
	}

	games := make([]AuditLogGameAffected, 0)
	err = tx.Select(&games, `
		SELECT a.*, w.name as white_name, b.name as black_name, g.played
		FROM audit_log_game_affected a
		JOIN games g ON a.game_ikey = g.ikey
		JOIN players w ON g.player_white = w.id
		JOIN players b ON g.player_black = b.id
		WHERE a.audit_log_id=$1
		ORDER BY is_main_target ASC, g.played ASC, g.ikey ASC;`, id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get audit log game affected"), err)
	}

	result := &DetailedAuditLog{
		AuditLog:  auditLog.AuditLog,
		AdminName: auditLog.AdminName,
		Players:   players,
		Games:     games,
	}

	return result, nil
}

type AuditLogUiFriendly struct {
	AuditLog
	AdminName string `db:"admin_name"`
}

func GetAuditLogsUiFriendly(db db.Db) ([]AuditLogUiFriendly, error) {
	logs := make([]AuditLogUiFriendly, 0)
	rows, err := db.GetSqlxDb().Queryx(`
		SELECT audit_logs.*, admin_users.name AS admin_name
		FROM audit_logs
		INNER JOIN admin_users ON admin_users.id = audit_logs.done_by
		ORDER BY created DESC;
		`)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get audit logs"), err)
	}

	for rows.Next() {
		var log AuditLogUiFriendly
		err = rows.StructScan(&log)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot scan audit log"), err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func GetAuditLogsUiFriendlyForPlayer(db db.Db, id uuid.UUID) ([]AuditLogUiFriendly, error) {
	logs := make([]AuditLogUiFriendly, 0)
	rows, err := db.GetSqlxDb().Queryx(`
		SELECT audit_logs.*, admin_users.name AS admin_name
		FROM audit_logs
		INNER JOIN admin_users ON admin_users.id = audit_logs.done_by
		INNER JOIN audit_log_player_affected ON audit_logs.id = audit_log_player_affected.audit_log_id
		INNER JOIN players ON players.id = audit_log_player_affected.player_id
		WHERE players.id = $1
		ORDER BY created DESC;
		`, id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get audit logs"), err)
	}

	for rows.Next() {
		var log AuditLogUiFriendly
		err = rows.StructScan(&log)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot scan audit log"), err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func GetAuditLogsUiFriendlyForGame(db db.Db, ikey int64) ([]AuditLogUiFriendly, error) {
	logs := make([]AuditLogUiFriendly, 0)
	rows, err := db.GetSqlxDb().Queryx(`
		SELECT audit_logs.*, admin_users.name AS admin_name
		FROM audit_logs
		INNER JOIN admin_users ON admin_users.id = audit_logs.done_by
		INNER JOIN audit_log_game_affected ON audit_logs.id = audit_log_game_affected.audit_log_id
		WHERE audit_log_game_affected.game_ikey = $1
		ORDER BY created DESC;
		`, ikey)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get audit logs"), err)
	}

	for rows.Next() {
		var log AuditLogUiFriendly
		err = rows.StructScan(&log)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot scan audit log"), err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func GetAuditLogsUiFriendlyForAdmin(db db.Db, id uuid.UUID) ([]AuditLogUiFriendly, error) {
	logs := make([]AuditLogUiFriendly, 0)
	rows, err := db.GetSqlxDb().Queryx(`
		SELECT audit_logs.*, admin_users.name AS admin_name
		FROM audit_logs
		INNER JOIN admin_users ON admin_users.id = audit_logs.done_by
		WHERE admin_users.id = $1
		ORDER BY created DESC;
		`, id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get audit logs"), err)
	}

	for rows.Next() {
		var log AuditLogUiFriendly
		err = rows.StructScan(&log)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot scan audit log"), err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}
