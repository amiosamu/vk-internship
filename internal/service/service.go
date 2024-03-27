package service

import (
	"context"
	"time"

	"github.com/amiosamu/vk-internship/internal/entity"
	"github.com/amiosamu/vk-internship/internal/repo"
	"github.com/google/uuid"
)



type AuthCreateUserInput struct {
	Name string
	Surname string
	Email string
	Password string
}

type AuthGenerateTokenInput struct {
	Email string
}


type Auth interface{
	RegisterUser(ctx context.Context, input AuthCreateUserInput) (uuid.UUID, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(token string) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
}


type Advertisement interface {
	CreateAdvertisement(advertisement entity.Advertisement) (uuid.UUID, error)
	GetAllAdvertisements(ctx context.Context, page, limit int, sortBy, sortOrder string, minPrice, maxPrice float64) ([]entity.Advertisement, error)
	GetAdvertisementByID(ctx context.Context, id uuid.UUID) (entity.Advertisement, error)
	GetAdvertisementsByUserID(ctx context.Context, id uuid.UUID) ([]entity.Advertisement, error)
}

type Services struct {
	Auth Auth
	Advertisement Advertisement
}

type ServiceDependencies struct {
	
	Repos *repo.Repos
	Signkey string
	TokenTTL time.Duration
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{

		Auth: NewAuthService(deps.Repos.User, deps.Signkey, deps.TokenTTL),
		Advertisement: NewAdvertisementService(deps.Repos.Advertisement),
	}
}