package repository

import (
	"database/sql"
	"time"
)

type Session struct {
	UUID            string    `json:"uuid"`
	ExpireDate      time.Time `json:"expireDate"`
	IdoniaAccountID uint32    `json:"idoniaAccountId"`
	IdoniaToken     string    `json:"idoniaToken"`
	IdoniaAPIKey    string    `json:"idoniaApiKey"`
	IdoniaAPISecret string    `json:"idoniaApiSecret"`
}

func InsertSession(db *sql.DB, session *Session) (err error) {
	_, err = db.Exec(
		"INSERT INTO session (uuid, expire_date, idonia_account_id, idonia_token, idonia_api_key, idonia_api_secret) VALUES ($1, $2, $3, $4, $5, $6)",
		session.UUID, session.ExpireDate, session.IdoniaAccountID, session.IdoniaToken, session.IdoniaAPIKey, session.IdoniaAPISecret,
	)
	return
}

func GetSession(db *sql.DB, uuid string) (session *Session, err error) {
	session = &Session{}
	res := db.QueryRow("SELECT uuid, expire_date, idonia_account_id, idonia_token, idonia_api_key, idonia_api_secret FROM session WHERE uuid = $1", uuid)
	err = res.Scan(&session.UUID, &session.ExpireDate, &session.IdoniaAccountID, &session.IdoniaToken, &session.IdoniaAPIKey, &session.IdoniaAPISecret)
	return
}

func EditSession(db *sql.DB, uuid string, session *Session) (err error) {
	_, err = db.Exec(
		"UPDATE session SET uuid = $1, expire_date = $2, idonia_account_id = $3, idonia_token = $4, idonia_api_key = $5, idonia_api_secret = $6 WHERE uuid = $7",
		session.UUID, session.ExpireDate, session.IdoniaAccountID, session.IdoniaToken, session.IdoniaAPIKey, session.IdoniaAPISecret, uuid,
	)
	return
}

func DeleteSession(db *sql.DB, uuid string) (err error) {
	_, err = db.Exec("DELETE FROM session WHERE uuid = $1", uuid)
	return
}

func DeleteAllSessions(db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM session")
	return
}
