package table

import (
	"database/sql"
	"errors"
	"fmt"

	"example.com/m/db"
)

// Represents any object that implements the Scan method, e.g. scannable objects.
type Scannable interface {
    Scan(dest ...any) error
}

// Interface meant for use by models that can be converted to and from sql rows.
type SqlCompliant[T any] interface {
    Create(scannable Scannable) (T, error)
    GetFields() (map[string]any, error)
}

// This implementation of table relies on FromScannable objects, meaning they can be created by scanning an sql row.
type postgres_table[T SqlCompliant[T]] struct {
    table_name string
    connection db.DB
}

// Selects the first instance of T by constructing a where clause. 'where' may be nil to select the first.
func (pt *postgres_table[T]) SelectFirst(where map[string]any) (T, error) {
    // Construct the where clause
    w := ""
    args := []any{}
    if where != nil {
        conditions := []string{}
        counter := 1
        for k, v := range where {
            conditions = append(conditions, fmt.Sprintf("%v=$%v", k, counter))
            args = append(args, v)
            counter++
        }

        for i, condition := range conditions {
            if i == 0 {
                w += fmt.Sprintf(" where %v", condition)
            } else {
                w += fmt.Sprintf(" and %v", condition)
            }
        }
    }

    // Construct query and execute it
    sql := fmt.Sprintf("select * from %v%v;", pt.table_name, w)
    row, err := pt.connection.QueryRow(sql, args...)

    if err != nil {
        var t T
        return t, err
    }

    // Create the object from a scannable object
    var t T
    t, err = t.Create(row)

    return t, err
}

// Selects all instance of T by constructing a where clause. 'where' may be nil to select all.
func (pt *postgres_table[T]) SelectAll(where map[string]any) ([]T, error) {
    // Construct the where clause
    w := ""
    args := []any{}
    if where != nil {
        conditions := []string{}
        counter := 1
        for k, v := range where {
            conditions = append(conditions, fmt.Sprintf("%v=$%v", k, counter))
            args = append(args, v)
            counter++
        }

        for i, condition := range conditions {
            if i == 0 {
                w += fmt.Sprintf(" where %v", condition)
            } else {
                w += fmt.Sprintf(" and %v", condition)
            }
        }
    }

    // Construct query and execute it
    sql := fmt.Sprintf("select * from %v%v;", pt.table_name, w)
    rows, err := pt.connection.QueryRows(sql, args...)

    ts := []T{}

    if err != nil {
        return ts, err
    }

    // Create the objects from a scannable object
    var t T
    for rows.Next() {
        t, err = t.Create(rows)

        if err != nil {
            return ts, err
        }

        ts = append(ts, t)
    }

    return ts, err
}

// Inserts the item T into the table. The GetFields method of T can also make restrictions by returning an error.
func (pt *postgres_table[T]) Insert(item T) error {
    kvps := item.GetFields()

    // Construct the clauses for the sql query
    nameClause := ""
    valueClause := ""
    values := []any{}
    counter := 1;
    for k, v := range kvps {
        if counter == 1 {
            nameClause += k
            valueClause += fmt.Sprintf("%v", counter)
        } else {
            nameClause += ", " + k
            valueClause += fmt.Sprintf(", %v", counter)
        }
        
        values = append(values, v)
        counter++
    }

    // Construct query and execute
    sql := fmt.Sprintf("insert into %v (%v) values (%v);", pt.table_name, nameClause, valueClause)
    res, err := pt.connection.Execute(sql, values...)

    if err != nil {
        return err
    }

    // If query affected 0 rows, return an error, otherwise, return the *maybe* error.
    affected, err := (*res).RowsAffected()

    if affected == 0 {
        return errors.New("Insert statement affected 0 rows!")
    }

    return err
}

// Deletes all instances matching the where statement constructed.
func (pt *postgres_table[T]) Delete(where map[string]any) (*sql.Result, error) {
    // Construct the where clause
    whereClause := ""
    args := []any{}
    if where != nil {
        conditions := []string{}
        counter := 1
        for k, v := range where {
            conditions = append(conditions, fmt.Sprintf("%v=$%v", k, counter))
            args = append(args, v)
            counter++
        }

        for i, condition := range conditions {
            if i == 0 {
                whereClause += fmt.Sprintf(" where %v", condition)
            } else {
                whereClause += fmt.Sprintf(" and %v", condition)
            }
        }
    }

    sql := fmt.Sprintf("delete from %v%v;", pt.table_name, whereClause)
    return pt.connection.Execute(sql, args...)
}

// Updates all the fields specified in set, for all instances matching the where clause.
func (pt *postgres_table[T]) Update(set map[string]any, where map[string]any) (*sql.Result, error) {
    // Construct the where clause
    whereClause := ""
    args := []any{}
    counter := 1
    if where != nil {
        conditions := []string{}
        for k, v := range where {
            conditions = append(conditions, fmt.Sprintf("%v=$%v", k, counter))
            args = append(args, v)
            counter++
        }

        for i, condition := range conditions {
            if i == 0 {
                whereClause += fmt.Sprintf(" where %v", condition)
            } else {
                whereClause += fmt.Sprintf(" and %v", condition)
            }
        }
    }

    setClause := ""
    flag := true
    for k, v := range set {
        if flag {
            setClause += fmt.Sprintf("%v=$%v", k, counter)
            flag = !flag
        } else {
            setClause += fmt.Sprintf(", %v=$%v", k, counter)
        }

        counter++
        args = append(args, v)
    }

    sql := fmt.Sprintf("update %v set %v%v;", pt.table_name, setClause, whereClause)
    return pt.connection.Execute(sql, args...)
}
