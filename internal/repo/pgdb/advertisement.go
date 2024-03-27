package pgdb

import (
	"context"
	"fmt"
	"time"

	"github.com/amiosamu/vk-internship/internal/entity"
	"github.com/amiosamu/vk-internship/internal/repo/repoerrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	maxPaginationLimit            = 10
	defaultPaginationLimit        = 10
	PriceSortType          string = "price"
	DateSortType           string = "date"
)

type AdvertisementRepo struct {
	DB *sqlx.DB
}

func (a *AdvertisementRepo) CreateAdvertisement(advertisement entity.Advertisement) (uuid.UUID, error) {
	q := `INSERT INTO advertisements (user_id, title, description, pictures, price) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{advertisement.UserID, advertisement.Title, advertisement.Description, pq.Array(&advertisement.Pictures), advertisement.Price}
	err := a.DB.QueryRowContext(ctx, q, args...).Scan(&advertisement.ID)
	if err != nil {
		return uuid.UUID{}, repoerrors.CannotCreate
	}
	return advertisement.ID, nil
}

func (a *AdvertisementRepo) GetAdvertisementByID(ctx context.Context, id uuid.UUID) (entity.Advertisement, error) {
	advertisement := entity.Advertisement{}
	q := `SELECT id, user_id, title, description, pictures, price, created_at FROM advertisements WHERE id = $1`
	err := a.DB.QueryRowContext(ctx, q, id).Scan(
		&advertisement.ID,
		&advertisement.UserID,
		&advertisement.Title,
		&advertisement.Description,
		pq.Array(&advertisement.Pictures),
		&advertisement.Price,
		&advertisement.CreatedAt,
	)
	if err != nil {
		return entity.Advertisement{}, repoerrors.ErrNotFound
	}
	return advertisement, nil
}

func (a *AdvertisementRepo) GetAllAdvertisements(ctx context.Context, page, limit int, sortBy, sortOrder string, minPrice, maxPrice float64) ([]entity.Advertisement, error) {
	var advertisements []entity.Advertisement
	q := fmt.Sprintf(`SELECT id, user_id, title, description, pictures, price, created_at 
	FROM advertisements 
	WHERE price >= $1 AND price <= $2 
	ORDER BY %s %s 
	LIMIT $3 OFFSET $4`, sortBy, sortOrder)

	offset := (page - 1) * limit

	rows, err := a.DB.QueryContext(ctx, q, minPrice, maxPrice, limit, offset)
	if err != nil {
		return advertisements, err
	}
	defer rows.Close()

	for rows.Next() {
		advertisement := entity.Advertisement{}
		err := rows.Scan(
			&advertisement.ID,
			&advertisement.UserID,
			&advertisement.Title,
			&advertisement.Description,
			pq.Array(&advertisement.Pictures),
			&advertisement.Price,
			&advertisement.CreatedAt,
		)
		if err != nil {
			return advertisements, err
		}

		advertisements = append(advertisements, advertisement)
	}

	if err := rows.Err(); err != nil {
		return advertisements, err
	}

	return advertisements, nil
}

func (a *AdvertisementRepo) GetAdvertisementsByUserID(ctx context.Context, id uuid.UUID) ([]entity.Advertisement, error) {
	var advertisements []entity.Advertisement
	q := `SELECT id, title, description, pictures, price, created_at FROM advertisements WHERE user_id = $1`
	rows, err := a.DB.QueryContext(ctx, q, id)
	if err != nil {
		return advertisements, err
	}
	defer rows.Close()

	for rows.Next() {
		advertisement := entity.Advertisement{}
		err := rows.Scan(
			&advertisement.ID,
			&advertisement.Title,
			&advertisement.Description,
			pq.Array(&advertisement.Pictures),
			&advertisement.Price,

			&advertisement.CreatedAt,
		)
		if err != nil {
			return advertisements, err
		}

		advertisements = append(advertisements, advertisement)
	}

	if err := rows.Err(); err != nil {
		return advertisements, err
	}

	return advertisements, nil
}

func NewAdvertisementRepo(pg *sqlx.DB) *AdvertisementRepo {
	return &AdvertisementRepo{pg}
}
