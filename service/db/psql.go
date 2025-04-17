package db

import (
	"database/sql"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

type DbAccess struct {
	db *sql.DB
}

func New(dsn string) (DbAccess, error) {
	db := DbAccess{}

	var err error
	db.db, err = sql.Open("pgx", dsn)
	sqlcreate := sqlbuilder.PostgreSQL.NewCreateTableBuilder()

	sql := sqlcreate.CreateTable("token").IfNotExists().
		Define("id", "TEXT", "NOT NULL", "PRIMARY KEY").
		Define("refresh", "TEXT", "NOT NULL").
		Define("guid", "TEXT", "NOT NULL").String()

	_, err = db.db.Exec(sql)
	return db, err
}

func (Db *DbAccess) Close() {
	err := Db.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// FIXME: Check кого? Что эта функция проверяет то?
func (Db *DbAccess) Check(id string) (string, error) {
	row := Db.db.QueryRow("SELECT refresh FROM token WHERE id = $1", id)

	var token string
	if err := row.Scan(&token); err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
	}

	return token, nil
}

// FIXME: Add кого?
func (Db *DbAccess) Add(guid, refresh string, id string) error {
	if _, err := Db.db.Exec("INSERT INTO token (guid, refresh, id) VALUES ($1,$2,$3)", guid, refresh, id); err != nil {
		return err
	}
	return nil
}

// FIXME: SetRefreshToken
func (Db *DbAccess) Refresh(newHashToken, id string) error {
	_, err := Db.db.Exec("UPDATE token SET refresh = $1 WHERE id = $2", newHashToken, id)
	return err
}

func (Db *DbAccess) GetEmail(guid string) string {
	row := Db.db.QueryRow("SELECT email FROM users WHERE guid = $1", guid)

	var email string
	if err := row.Scan(&email); err != nil {
		if err != sql.ErrNoRows {
			return ""
		}
	}

	return email
}

// FIXME: Delete кого? По названию метода неясно.
func (Db *DbAccess) Delete(id string) error {
	_, err := Db.db.Exec("DELETE FROM token WHERE id = $1", id)
	return err
}
