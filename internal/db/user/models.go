// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package user

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Bio       sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
