package db

import "database/sql"

// Interface for simple sql database connection
type DB interface {
    QueryRow(query string, args ...any) (*sql.Row, error)
    QueryRows(query string, args ...any) (*sql.Rows, error)
    Execute(query string, args ...any) (*sql.Result, error)
    AddPreProcess(f func(query string, args ...any))
    AddPostProcess(f func(query string, args ...any))
}
