package table

import "database/sql"

// This interface provides functions to alter a table in sql,
// with the most basic commands covered.
type Table[T any] interface {
    // Select only get the first matching instance of T.
    SelectFirst(where map[string]any) (T, error)

    // Select all matching instances of T.
    SelectAll(where map[string]any) ([]T, error)

    // Insert T into the table
    Insert(item T) error

    // Delete all matching instances of T.
    Delete(where map[string]any) (*sql.Result, error)

    // Update the fields in set with the mapped value for all matching instances of T.
    Update(set map[string]any, where map[string]any) (*sql.Result, error)
}
