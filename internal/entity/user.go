package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Surname   string    `db:"surname"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
