package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hassamk122/restapi_golang/internal/auth"
	"github.com/hassamk122/restapi_golang/internal/repo"
	"github.com/hassamk122/restapi_golang/internal/store"
	"github.com/hassamk122/restapi_golang/internal/utils"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	DB       *sql.DB
	UserRepo repo.UserRepo
	Redis    *redis.Client
}

func NewUserService(db *sql.DB, userRepo repo.UserRepo, redisClient *redis.Client) *UserService {
	return &UserService{
		DB:       db,
		UserRepo: userRepo,
		Redis:    redisClient,
	}
}

func (s *UserService) Register(ctx context.Context, username, email, password string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := store.New(tx)
	repo := repo.NewUserRepo(qtx)

	_, err = repo.GetUserByEmail(ctx, email)
	if err == nil {
		return ErrEmailTaken
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = repo.CreateUser(ctx, store.CreateUserParams{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	})

	return tx.Commit()
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.UserRepo.GetUserByEmailIncludingPassword(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if !utils.ComparePassword(user.Password, password) {
		return "", ErrInvalidCredentials
	}
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	return auth.GenerateJWT(int(user.ID), user.Email, jwtKey)
}

func (s *UserService) GetUserProfile(ctx context.Context, userId int32) (*store.GetUserRow, error) {
	cacheKey := fmt.Sprintf("user:%d", userId)
	if cached, err := s.Redis.Get(ctx, cacheKey).Result(); err == nil {
		var user store.GetUserRow
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			log.Println("returned from redis cache")
			return &user, nil
		}
	}

	user, err := s.UserRepo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	if err := s.Redis.Set(ctx, cacheKey, userJson, 5*time.Minute).Err(); err != nil {
		log.Println("redis set failed:", err)
	}

	return &user, nil
}
