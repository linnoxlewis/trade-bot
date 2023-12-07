package service

import (
	"context"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

var errUserNotFound = "userNotFound"

type ApiKeyRepo interface {
	AddApiKeys(ctx context.Context, keys *domain.ApiKeys) (err error)
	DeleteApiKey(ctx context.Context, userId int64, exchange string) (err error)
	ClearApiKey(ctx context.Context, userId int64) (err error)
	GetApiKeysByUserIdAndExchange(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error)
	GetApiKeysByUserId(ctx context.Context, userId int64) (apiKeys []*domain.ApiKeys, err error)
}

type ApiKeysService struct {
	cfg        *config.Config
	apiKeyRepo ApiKeyRepo
	userRepo   UserRepo
	i18n       *i18n.I18n
	logger     *log.Logger
}

func NewApiKeysService(cfg *config.Config,
	apiKeyRepo ApiKeyRepo,
	userRepo UserRepo,
	i18n *i18n.I18n,
	logger *log.Logger) *ApiKeysService {
	return &ApiKeysService{cfg,
		apiKeyRepo,
		userRepo,
		i18n,
		logger}
}

func (a *ApiKeysService) AddApiKeys(ctx context.Context, apiKeys *dto.ApiKeys) error {
	if !a.userRepo.ExistUser(ctx, apiKeys.UserId) {
		return errors.BadRequestError(a.i18n.T(errUserNotFound, nil, "ru"))
	}

	privateKey, err := helper.DecryptClientMessage(apiKeys.PrivKey, a.cfg.GetClientSecretKey())
	if err != nil {
		a.logger.ErrorLog.Println("err decrypt priv key: ", err)

		return errors.InternalServerError(err)
	}
	encryptPrivKey, err := helper.EncryptMessage(privateKey, a.cfg.GetApiSecret())
	if err != nil {
		a.logger.ErrorLog.Println("err encrypt private key: ", err)

		return errors.InternalServerError(err)
	}

	var passphrase string
	if !apiKeys.EmptyPassPhrase() {
		passphrase, err = helper.DecryptClientMessage(apiKeys.PassPhrase, a.cfg.GetClientSecretKey())
		if err != nil {
			a.logger.ErrorLog.Println("err decrypt passphrase key: ", err)

			return errors.InternalServerError(err)
		}

		passphrase, err = helper.EncryptMessage(passphrase, a.cfg.GetApiSecret())
		if err != nil {
			a.logger.ErrorLog.Println("err encrypt private key: ", err)

			return errors.InternalServerError(err)
		}
	}

	userApiKey := domain.NewApiKeys(apiKeys.UserId, apiKeys.Exchange, apiKeys.PubKey, encryptPrivKey, passphrase)
	if err = a.apiKeyRepo.AddApiKeys(ctx, userApiKey); err != nil {
		a.logger.ErrorLog.Println("err add api keys: ", err)

		return errors.InternalServerError(err)
	}

	return nil
}

func (a *ApiKeysService) DeleteApiKey(ctx context.Context, userId int64, exchange string) error {
	if err := a.apiKeyRepo.DeleteApiKey(ctx, userId, exchange); err != nil {
		a.logger.ErrorLog.Println("err delete api key", err)

		return errors.InternalServerError(err)
	}

	return nil
}

func (a *ApiKeysService) ClearApiKey(ctx context.Context, userId int64, exchange string) error {
	if err := a.apiKeyRepo.ClearApiKey(ctx, userId); err != nil {
		a.logger.ErrorLog.Println("err clear api key", err)

		return errors.InternalServerError(err)
	}

	return nil
}

func (a *ApiKeysService) GetApiKeyByExchangeAndId(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error) {
	apiKeys, err := a.apiKeyRepo.GetApiKeysByUserIdAndExchange(ctx, userId, exchange)
	if err != nil {
		a.logger.ErrorLog.Println("err get  api key", err)

		return nil, errors.InternalServerError(err)
	}

	return apiKeys, nil
}
