package db

import (
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
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

	sqlcreate = sqlbuilder.PostgreSQL.NewCreateTableBuilder()
	sql = sqlcreate.CreateTable("usersMail").IfNotExists().
		Define("guid", "TEXT", "NOT NULL", "PRIMARY KEY").
		Define("email", "TEXT", "NOT NULL").String()

	_, err = db.db.Exec(sql)

	return db, err
}

func (Db *DbAccess) Close() {
	err := Db.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (Db *DbAccess) GetRefreshToken(id string) (string, error) {
	row := Db.db.QueryRow("SELECT refresh FROM token WHERE id = $1", id)

	var token string
	if err := row.Scan(&token); err != nil {
		if err != sql.ErrNoRows {
			return "", fmt.Errorf("Error getting refresh token: %v", err)
		}
	}

	return token, nil
}

func (Db *DbAccess) AddToken(guid, id string, refreshToken []byte) error {

	RefreshBcrypt, err := bcrypt.GenerateFromPassword(refreshToken, bcrypt.DefaultCost)
	if err != nil {
		fmt.Errorf("Generation bcrypt failed: %s", err)
	}

	if _, err = Db.db.Exec("INSERT INTO token (guid, refresh, id) VALUES ($1,$2,$3)", guid, RefreshBcrypt, id); err != nil {
		return fmt.Errorf("Error inserting token: %s", err)
	}
	return nil
}

func (Db *DbAccess) SetRefreshToken(id string, refreshToken []byte) error {
	NewRefreshBcrypt, err := bcrypt.GenerateFromPassword(refreshToken, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Generation bcrypt failed: %s", err)
	}

	_, err = Db.db.Exec("UPDATE token SET refresh = $1 WHERE id = $2", string(NewRefreshBcrypt), id)
	if err != nil {

	}
	return err
}

func (Db *DbAccess) GetEmail(guid string) string {
	row := Db.db.QueryRow("SELECT email FROM usersMail WHERE guid = $1", guid)

	var email string
	if err := row.Scan(&email); err != nil {
		if err != sql.ErrNoRows {
			return ""
		}
	}

	return email
}

func (Db *DbAccess) DeleteToken(id string) error {
	_, err := Db.db.Exec("DELETE FROM token WHERE id = $1", id)
	return err
}
