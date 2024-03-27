package service

import (
	"context"
	"errors"

	"github.com/amiosamu/vk-internship/internal/entity"
	"github.com/amiosamu/vk-internship/internal/repo"
	"github.com/amiosamu/vk-internship/internal/repo/repoerrors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AdvertisementService struct {
	advertisementRepo repo.Advertisement
}

func (s *AdvertisementService) CreateAdvertisement(advertisement entity.Advertisement) (uuid.UUID, error) {
	id, err := s.advertisementRepo.CreateAdvertisement(advertisement)
	if err != nil {
		log.Fatalf("error: cannot create advertisement: %v\n", err.Error())
		return uuid.UUID{}, err
	}
	return id, nil
}

func (s *AdvertisementService) GetAdvertisementByID(ctx context.Context, id uuid.UUID) (entity.Advertisement, error) {
	advertisement, err := s.advertisementRepo.GetAdvertisementByID(ctx, id)
	if err != nil {

		if errors.Is(err, repoerrors.ErrNotFound) {
			return entity.Advertisement{}, ErrAdvertisementNotFound
		}
		log.Fatalf("Failed to get advertisement with ID %d: %s\n", id, err.Error())
	}
	return advertisement, nil
}

func (s *AdvertisementService) GetAllAdvertisements(ctx context.Context,  page, limit int, sortBy, sortOrder string, minPrice, maxPrice float64) ([]entity.Advertisement, error) {
	advertisements, err := s.advertisementRepo.GetAllAdvertisements(ctx, page, limit, sortBy, sortOrder, minPrice, maxPrice)
	if err != nil {
		log.Printf("Failed to get all advertisements: %v", err)
		return nil, ErrCannotGetAdvertisement
	}
	return advertisements, nil

}

func (s *AdvertisementService) GetAdvertisementsByUserID(ctx context.Context, id uuid.UUID) ([]entity.Advertisement, error) {
	advertisements, err := s.advertisementRepo.GetAdvertisementsByUserID(ctx, id)
	if err != nil  {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return []entity.Advertisement{}, ErrAdvertisementNotFound
		}
		return []entity.Advertisement{}, ErrCannotGetAdvertisement
	}
	return advertisements, nil
}

func NewAdvertisementService(advertisementRepo repo.Advertisement) *AdvertisementService {
	return &AdvertisementService{
		advertisementRepo: advertisementRepo,
	}
}
