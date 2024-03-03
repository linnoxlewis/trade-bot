package service

import (
	"context"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

var userAlreadyExist = "userAlreadyExist"

type UserRepo interface {
	CreateUser(ctx context.Context, user *domain.User) (err error)
	ExistUser(ctx context.Context, id int64) (exist bool)
	IsAdmin(ctx context.Context, id int64) (bool, error)
	GetAdminIds(ctx context.Context) ([]int, error)
}

type UserService struct {
	cfg      *config.Config
	userRepo UserRepo
	i18n     *i18n.I18n
	logger   *log.Logger
}

func NewUserService(cfg *config.Config,
	userRepo UserRepo,
	i18n *i18n.I18n,
	logger *log.Logger) *UserService {
	return &UserService{cfg,
		userRepo,
		i18n,
		logger}
}

func (u *UserService) CreateUser(ctx context.Context, username string, tgId int64) (err error) {
	if u.userRepo.ExistUser(ctx, tgId) {
		return errors.BadRequestError(u.i18n.T(userAlreadyExist, nil, helper.GetDefaultLg()))
	}

	if err = u.userRepo.CreateUser(ctx, &domain.User{ID: tgId, Username: username}); err != nil {
		u.logger.ErrorLog.Println("err can`t create user: ", err)

		return errors.InternalServerError(err)
	}

	return
}

func (u *UserService) IsAdmin(ctx context.Context, tgId int64) bool {
	if !u.userRepo.ExistUser(ctx, tgId) {
		return false
	}

	isAdmin, err := u.userRepo.IsAdmin(ctx, tgId)
	if err != nil {
		u.logger.ErrorLog.Println("err can`get admin user: ", err)

		return false
	}

	return isAdmin
}

func (u *UserService) GetAdmins(ctx context.Context) ([]int, error) {
	admins, err := u.userRepo.GetAdminIds(ctx)
	if err != nil {
		u.logger.ErrorLog.Println("err can`get admin ids: ", err)

		return nil, err
	}

	return admins, nil
}
