package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

type UserStorageInterface interface {
	GetUserById(ctx context.Context, id uuid.UUID) (user *domain.User, err error)
	CreateUser(ctx context.Context, user *domain.User) (err error)
}

type Userer interface {
	GetUserById(ctx context.Context, id uuid.UUID) (user *domain.User, err error)
	CreateUser(ctx context.Context, username string, tgId int64) (err error)
}

type UserService struct {
	cfg      *config.Config
	userRepo UserStorageInterface
	logger   log.Logger
}

func NewUserService(cfg *config.Config,
	userRepo UserStorageInterface,
	logger log.Logger) *UserService {
	return &UserService{cfg, userRepo, logger}
}

func (u *UserService) CreateUser(ctx context.Context, username string, tgId int64) (err error) {
	user := &domain.User{
		Id:       uuid.New(),
		TgId:     tgId,
		Username: username,
	}
	if err = u.userRepo.CreateUser(ctx, user); err != nil {
		u.logger.ErrorLog.Println("err can`t create user: ", err)

		return
	}

	return
}

func (u *UserService) GetUserById(ctx context.Context, id uuid.UUID) (user *domain.User, err error) {
	return u.userRepo.GetUserById(ctx, id)
}
