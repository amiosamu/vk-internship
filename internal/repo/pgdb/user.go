package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/amiosamu/vk-internship/internal/entity"
	"github.com/amiosamu/vk-internship/internal/repo/repoerrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo struct {
	DB *sqlx.DB
}

func (u *UserRepo) CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error) {
	userID := uuid.New()
	q := `INSERT INTO users (id, name, surname, email, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var insertedUserID string
	err := u.DB.QueryRowContext(ctx, q, userID.String(), user.Name, user.Surname, user.Email, user.Password).Scan(&insertedUserID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			log.Printf("could not create user with email: %v. reason: %v\n", user.Email, pgErr.Message)
			return uuid.Nil, repoerrors.ErrAlreadyExists
		}
		log.Printf("failed to create user: %v\n", err)
		return uuid.Nil, repoerrors.CannotCreate
	}

	insertedUUID, err := uuid.Parse(insertedUserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse inserted UUID: %w", err)
	}

	return insertedUUID, nil
}

func (u *UserRepo) getUserByQuery(ctx context.Context, query string, args ...interface{}) (entity.User, error) {
	user := entity.User{}

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Surname, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, repoerrors.ErrNotFound
		}
		return entity.User{}, repoerrors.CannotCreate
	}
	return user, nil
}

func (u *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	q := `SELECT id, name, surname, email FROM users WHERE id = $1`
	return u.getUserByQuery(ctx, q, id)
}

func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	q := `SELECT id, name, surname, email FROM users WHERE email = $1`
	return u.getUserByQuery(ctx, q, email)
}

func NewUserRepo(pg *sqlx.DB) *UserRepo {
	return &UserRepo{DB: pg}
}
