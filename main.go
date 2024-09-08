package main

import (
	"fmt"
	"os"
	"time"

	"example.com/m/db"
	"example.com/m/models"
	"example.com/m/table"
	_ "github.com/lib/pq"
)

func main() {
    db := getConnection()
    db.AddPreProcess(func(query string, args ...any) {
        fmt.Printf("Executing: %v\n", query)
        fmt.Printf("Args: %v\n", args)
    })

    dropAll(db)
    createUserTable(db)
    createOffices(db)
    createUsers(db)
    selectAUser(db)
    updateAUser(db)
}

func dropAll(connection db.DB) {
    _, err := connection.Execute("drop table if exists users;")
    _, err = connection.Execute("drop table if exists offices;")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Couldn't drop databases: %v\n", err)
        os.Exit(1)
    }
}

func createUserTable(connection db.DB) {
    _, err := connection.Execute("create table if not exists users (id serial primary key, first_name text, last_name text, birth_date bigint, office_id bigint);")
    if err != nil {
        fmt.Fprintf(os.Stderr, "User table was not created: %v\n", err)
        os.Exit(1)
    }
}

func createOffices(connection db.DB) {
    _, err := connection.Execute("create table if not exists offices (id serial primary key, name text, capacity int);")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Offices not created: %v\n", err)
        os.Exit(1)
    }
}

func getConnection() db.DB {
    db, err := db.CreatePostgresDB("mads", "mads", "mads", 666, false)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Connection was not established: %v\n", err)
        os.Exit(1)
    }
    return db
}

func createUsers(connection db.DB) {
    userTable := table.CreatePostgresTable[*models.User]("users", connection)
    
    for i := range 100 {
        user := models.CreateUser("mads hvid", "byriel", time.Now().Unix(), int64(i))
        err := userTable.Insert(user)

        if err != nil {
            fmt.Fprintf(os.Stderr, "userTable.Insert error: %v\n", err)
        }
    }
}

func insertBadUser(connection db.DB) {
    userTable := table.CreatePostgresTable[*models.User]("users", connection)
    user := models.CreateUser("mads hvid", "byriel", -10000000000, int64(0))
    err := userTable.Insert(user)

    if err != nil {
        fmt.Fprintf(os.Stderr, "userTable.Insert error: %v\n", err)
        return 
    }
}

func selectAUser(connection db.DB) {
    userTable := table.CreatePostgresTable[*models.User]("users", connection)
    user, err := userTable.SelectFirst(map[string]any{"first_name": "mads hvid"})
    if err != nil {
        fmt.Fprintf(os.Stderr, "userTable.SelectFirst error: %v\n", err)
        return 
    }
    fmt.Printf("Selected user: %+v\n", user)
}

func updateAUser(connection db.DB) {
    // Pick some random user
    userTable := table.CreatePostgresTable[*models.User]("users", connection)
    user, err := userTable.SelectFirst(map[string]any{"first_name": "mads hvid"})
    if err != nil {
        fmt.Fprintf(os.Stderr, "userTable.SelectFirst error: %v\n", err)
        return 
    }

    _, err = userTable.Update(map[string]any{"first_name": "find", "last_name": "holger"}, map[string]any{"id": user.GetId()})

    if err != nil {
        fmt.Fprintf(os.Stderr, "userTable.Update error: %v\n", err)
        return 
    }
}
