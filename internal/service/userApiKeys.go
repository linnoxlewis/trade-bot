package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

type ApiKeyRepo interface {
	AddApiKeys(ctx context.Context, keys *domain.ApiKeys) (err error)
	DeleteApiKey(ctx context.Context, userId uuid.UUID, exchange string) (err error)
	ClearApiKey(ctx context.Context, userId uuid.UUID) (err error)

	GetApiKeysByUserIdAndExchange(ctx context.Context, userId uuid.UUID, exchange string) (apiKeys *domain.ApiKeys, err error)
	GetByApiKeyByTgUserId(ctx context.Context, tgId int64) (apiKeys []*domain.ApiKeys, err error)
	GetByApiKeyByTgUserIdAndExchange(ctx context.Context, tgId int, exchange string) (apiKeys *domain.ApiKeys, err error)
	GetApiKeysByUserId(ctx context.Context, userId uuid.UUID) (apiKeys []*domain.ApiKeys, err error)
}

type ApiKeysService struct {
	cfg         *config.Config
	apiKeyRepo  ApiKeyRepo
	userService Userer
	logger      *log.Logger
}

func NewApiKeysService(cfg *config.Config,
	apiKeyRepo ApiKeyRepo,
	userService Userer,
	logger *log.Logger) *ApiKeysService {
	return &ApiKeysService{cfg, apiKeyRepo, userService, logger}
}

func (a *ApiKeysService) AddApiKeys(ctx context.Context, userId uuid.UUID, apiKeys *dto.ApiKeys) error {
	user, err := a.userService.GetUserById(ctx, userId)
	if err != nil {
		a.logger.ErrorLog.Println("can`t get user: ", err)

		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	privateKey, err := helper.DecryptClientMessage(apiKeys.PrivKey, a.cfg.GetClientSecretKey())
	if err != nil {
		a.logger.ErrorLog.Println("err decrypt priv key: ", err)

		return err
	}
	encryptPrivKey, err := helper.EncryptMessage(privateKey, a.cfg.GetSecretKey())
	if err != nil {
		a.logger.ErrorLog.Println("err encrypt private key: ", err)

		return err
	}

	var passphrase string
	if !apiKeys.EmptyPassPhrase() {
		passphrase, err = helper.DecryptClientMessage(apiKeys.PassPhrase, a.cfg.GetClientSecretKey())
		if err != nil {
			a.logger.ErrorLog.Println("err decrypt passphrase key: ", err)

			return err
		}

		passphrase, err = helper.EncryptMessage(passphrase, a.cfg.GetSecretKey())
		if err != nil {
			a.logger.ErrorLog.Println("err encrypt private key: ", err)

			return err
		}
	}

	userApiKey := domain.NewApiKeys(userId, apiKeys.Exchange, apiKeys.PubKey, encryptPrivKey, passphrase)
	if err = a.apiKeyRepo.AddApiKeys(ctx, userApiKey); err != nil {
		a.logger.ErrorLog.Println("err add api keys: ", err)

		return err
	}

	return nil
}

func (a *ApiKeysService) DeleteApiKey(ctx context.Context, userId uuid.UUID, exchange string) error {
	return a.apiKeyRepo.DeleteApiKey(ctx, userId, exchange)
}

func (a *ApiKeysService) ClearApiKey(ctx context.Context, userId uuid.UUID, exchange string) error {
	return a.apiKeyRepo.ClearApiKey(ctx, userId)
}

func (a *ApiKeysService) GetApiKeyByExchangeAndTgId(ctx context.Context, tgId int64, exchange string) (*domain.ApiKeys, error) {
	return a.apiKeyRepo.GetByApiKeyByTgUserIdAndExchange(ctx, int(tgId), exchange)
}
