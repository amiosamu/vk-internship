package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/amiosamu/vk-internship/internal/entity"
	"github.com/amiosamu/vk-internship/internal/repo"
	"github.com/amiosamu/vk-internship/internal/repo/repoerrors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID
}

type AuthService struct {
	userRepo repo.User
	signKey  string
	tokenTTL time.Duration
}

func NewAuthService(userRepo repo.User, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		signKey:  signKey,
		tokenTTL: tokenTTL,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, input AuthCreateUserInput) (uuid.UUID, error) {
	user := entity.User{
		Name:     input.Name,
		Surname:  input.Email,
		Email:    input.Email,
		Password: input.Password,
	}

	userID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		if err == repoerrors.ErrAlreadyExists {
			return uuid.UUID{}, err
		}
		log.Errorf("AuthService.RegisterUser - c.UserRepo.CreateUser: %v", err)
	}

	return userID, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return "", ErrUserNotFound
		}
		log.Errorf("AuthService.GenerateToken: cannot get user: %v", err)
		return "", ErrCannotGetUser
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  int64(time.Now().Unix()),
		},
		UserID: user.ID,
	})

	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		log.Errorf("AuthService.GenerateToken: cannot sign token: %v", err)
		return "", ErrCannotSignToken
	}
	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return uuid.UUID{}, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return uuid.UUID{}, ErrCannotParseToken
	}

	return claims.UserID, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			return entity.User{}, ErrUserNotFound
		}
		return entity.User{}, ErrCannotGetUser
	}

	return user, nil
}
