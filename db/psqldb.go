package db

import (
	"database/sql"
	"fmt"
)

type psqlDB struct {
    connection *sql.DB
    preProcesses []func(string, ...any)
    postProcesses []func(string, ...any)
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

    pDB := &psqlDB {connection: db, preProcesses: make([]func(string, ...any), 0), postProcesses: make([]func(string, ...any), 0)}

    return pDB, err
}

func sslModeConvert(sslmode bool) string {
    if sslmode {
        return "enable"
    }
    return "disable"
}

func (p *psqlDB) QueryRow(query string, args ...any) (*sql.Row, error) {
    for _, f := range p.preProcesses {
        f(query, args...)
    }
    for _, f := range p.postProcesses {
        defer f(query, args...)
    }

    return p.connection.QueryRow(query, args...), nil
}

func (p *psqlDB) QueryRows(query string, args ...any) (*sql.Rows, error) {
    for _, f := range p.preProcesses {
        f(query, args...)
    }
    for _, f := range p.postProcesses {
        defer f(query, args...)
    }

    return p.connection.Query(query, args...)
}

func (p *psqlDB) Execute(query string, args ...any) (*sql.Result, error) {
    for _, f := range p.preProcesses {
        f(query, args...)
    }
    for _, f := range p.postProcesses {
        defer f(query, args...)
    }

    res, err := p.connection.Exec(query, args...)
    return &res, err
}

func (p *psqlDB) AddPreProcess(f func(query string, args ...any)) {
    p.preProcesses = append(p.preProcesses, f)
}

func (p *psqlDB) AddPostProcess(f func(query string, args ...any)) {
    p.postProcesses = append(p.preProcesses, f)
}
