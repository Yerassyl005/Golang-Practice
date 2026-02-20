package modules

import "time"

type User struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Email     string     `db:"email" json:"email"`
	Age       int        `db:"age" json:"age"`
	City      string     `db:"city" json:"city"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}