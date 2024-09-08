package db

import (
	"database/sql"
	"fmt"
)

type psqlDB struct {
    connection *sql.DB
}

// Creates a connection to a postgesql database with specified user, passwor
func CreatePostgresDB(user string, password string, dbname string, port int, sslmode bool) (DB, error) {
    db, err := sql.Open(
        "postgres",
        fmt.Sprintf(
            "user=%v password=%v dbname=%v port=%v sslmode=%v",
            user,
            password,
            dbname,
            port,
            sslModeConvert(sslmode),
        ),
    )

    pDB := &psqlDB {connection: db}

    return pDB, err
}

func sslModeConvert(sslmode bool) string {
    if sslmode {
        return "enable"
    }
    return "disable"
}

func (p *psqlDB) QueryRow(query string, args ...any) (*sql.Row, error) {
    return p.connection.QueryRow(query, args...), nil
}

func (p *psqlDB) QueryRows(query string, args ...any) (*sql.Rows, error) {
    return p.connection.Query(query, args...)
}

func (p *psqlDB) Execute(query string, args ...any) (*sql.Result, error) {
    res, err := p.connection.Exec(query, args...)
    return &res, err
}
