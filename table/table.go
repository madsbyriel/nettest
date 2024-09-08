package table

import "database/sql"

type Table[T any] interface {
    SelectFirst(where map[string]any) (T, error)
    SelectAll(where map[string]any) ([]T, error)
    Insert(item T) error
    Delete(where map[string]any) (*sql.Result, error)
    Update(set map[string]any, where map[string]any) (*sql.Result, error)
}
