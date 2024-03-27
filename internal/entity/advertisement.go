package entity

import (
	"time"

	"github.com/google/uuid"
)

type Advertisement struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Pictures    []string  `db:"pictures"`
	Address     string    `db:"address"`
	Price       float64   `db:"price"`
	CreatedAt   time.Time `db:"created_at"`
}
