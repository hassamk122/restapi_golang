package service

import (
	"context"
	"database/sql"
	"os"

	"github.com/hassamk122/restapi_golang/internal/auth"
	"github.com/hassamk122/restapi_golang/internal/repo"
	"github.com/hassamk122/restapi_golang/internal/store"
	"github.com/hassamk122/restapi_golang/internal/utils"
)

type UserService struct {
	DB       *sql.DB
	UserRepo repo.UserRepo
}

func NewUserService(db *sql.DB, userRepo repo.UserRepo) *UserService {
	return &UserService{
		DB:       db,
		UserRepo: userRepo,
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
	user, err := s.UserRepo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}
