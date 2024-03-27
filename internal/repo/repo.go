package repo

import (
	"context"

	"github.com/amiosamu/vk-internship/internal/entity"
	"github.com/amiosamu/vk-internship/internal/repo/pgdb"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
}

type Advertisement interface {
	CreateAdvertisement(advertisement entity.Advertisement) (uuid.UUID, error)
	GetAllAdvertisements(ctx context.Context, page, limit int, sortBy, sortOrder string, minPrice, maxPrice float64) ([]entity.Advertisement, error)
	GetAdvertisementByID(ctx context.Context, id uuid.UUID) (entity.Advertisement, error)
	GetAdvertisementsByUserID(ctx context.Context, id uuid.UUID) ([]entity.Advertisement, error)
}

type Repos struct {
	Advertisement
	User
}

func NewRepos(pg *sqlx.DB) *Repos {
	return &Repos{
		User:          pgdb.NewUserRepo(pg),
		Advertisement: pgdb.NewAdvertisementRepo(pg),
	}
}
