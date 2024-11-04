package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

// Ошибки доменного уровня
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserCreation      = errors.New("failed to create user")
	ErrAdminCheck        = errors.New("failed to check if user is admin")
	ErrFetchAdmins       = errors.New("failed to fetch admin IDs")
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *domain.User) error
	ExistUser(ctx context.Context, id int64) bool
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
	return &UserService{
		cfg:      cfg,
		userRepo: userRepo,
		i18n:     i18n,
		logger:   logger,
	}
}

// CreateUser создает нового пользователя
func (u *UserService) CreateUser(ctx context.Context, username string, tgId int64) error {
	if u.userRepo.ExistUser(ctx, tgId) {
		u.logger.ErrorLog.Printf("User with tgId %d already exists", tgId)
		return fmt.Errorf("%w: tgId %d", ErrUserAlreadyExists, tgId)
	}

	user := &domain.User{ID: tgId, Username: username}
	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		u.logger.ErrorLog.Printf("Failed to create user with tgId %d: %v", tgId, err)
		return fmt.Errorf("%w: %v", ErrUserCreation, err)
	}

	u.logger.InfoLog.Printf("User with tgId %d created successfully", tgId)
	return nil
}

// IsAdmin проверяет, является ли пользователь администратором
func (u *UserService) IsAdmin(ctx context.Context, tgId int64) bool {
	if !u.userRepo.ExistUser(ctx, tgId) {
		u.logger.ErrorLog.Printf("User with tgId %d does not exist", tgId)
		return false
	}

	isAdmin, err := u.userRepo.IsAdmin(ctx, tgId)
	if err != nil {
		u.logger.ErrorLog.Printf("Failed to check if user with tgId %d is admin: %v", tgId, err)
		return false
	}

	return isAdmin
}

// GetAdmins возвращает список ID администраторов
func (u *UserService) GetAdmins(ctx context.Context) ([]int, error) {
	admins, err := u.userRepo.GetAdminIds(ctx)
	if err != nil {
		u.logger.ErrorLog.Printf("Failed to fetch admin IDs: %v", err)
		return nil, fmt.Errorf("%w: %v", ErrFetchAdmins, err)
	}

	u.logger.InfoLog.Printf("Fetched %d admin IDs successfully", len(admins))
	return admins, nil
}
