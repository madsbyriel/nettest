package models

import (
	"errors"
	"strings"
	"time"

	"example.com/m/table"
)

type User struct {
    id int64
    first_name string
    last_name string
    birth_date int64
    office_id int64
}

func CreateUser(first_name, last_name string, birth_date, office_id int64) *User {
    return &User{0, first_name, last_name, birth_date, office_id}
}

// Creates a new user from a scannable object, most likely an sql row.
func (u *User) Create(scannable table.Scannable) (*User, error) {
    user := &User{}
    err := scannable.Scan(&user.id, &user.first_name, &user.last_name, &user.birth_date, &user.office_id)
    return user, err
}

// Function used to insert new users into the user table. Let it crash whenever a field is illegal.
// Id's are not supplied by this method as the table should generate those by itself.
func (u *User) GetFields() (map[string]any, error) {
    // Could use some extra validation
    if len(u.first_name) == 0 {
        return nil, errors.New("First name of this user is empty!")
    }

    if strings.Contains(u.last_name, " ") {
        return nil, errors.New("White space in last names is not allowed, middle names must be contained in first_name field.")
    }

    // People over the age of 120 are now allowed (because they are most likely vampires)
    yearDelta := time.Now().Year() - time.Unix(u.birth_date, 0).Year()
    if yearDelta > 120 {
        return nil, errors.New("Illegal birth date, this person is over 120 years old!")
    }
    if yearDelta < 0 {
        return nil, errors.New("Illegal birth date, this person has yet to be born!")
    }

    fields := map[string]any{
        "first_name": u.first_name,
        "last_name": u.last_name, 
        "birth_date": u.birth_date,
        "office_id": u.office_id,
    }

    return fields, nil
}
