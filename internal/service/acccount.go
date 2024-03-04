package service

import (
	"context"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/internal/pkg/exchanger"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

type AccountService struct {
	cfg        *config.Config
	apiKeyRepo ApiKeyRepo
	exchanger  exchanger.Exchanger
	i18n       *i18n.I18n
	logger     *log.Logger
}

func NewAccountService(cfg *config.Config,
	apiKeyRepo ApiKeyRepo,
	exchanger exchanger.Exchanger,
	i18n *i18n.I18n,
	logger *log.Logger) *AccountService {
	return &AccountService{
		cfg:        cfg,
		apiKeyRepo: apiKeyRepo,
		exchanger:  exchanger,
		i18n:       i18n,
		logger:     logger}
}

func (a *AccountService) GetBalance(ctx context.Context, userId int64, exchange string) (balance domain.Balance, err error) {
	keys, err := a.getApiKeys(ctx, userId, exchange)
	if err != nil {
		return nil, err
	}
	result, err := a.exchanger.Balance(keys, exchange)
	if err != nil {
		return nil, errors.BadRequestError(err.Error())
	}

	return result, nil
}

func (a *AccountService) getApiKeys(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error) {
	keys, err := a.apiKeyRepo.GetApiKeysByUserIdAndExchange(ctx, userId, exchange)
	if err != nil {
		a.logger.ErrorLog.Println("err get api keys: ", err)

		return nil, errors.InternalServerError(err)
	}
	if keys == nil {
		return nil, errors.BadRequestError(a.i18n.T(errApiKeysNotFound, nil, "ru"))
	}
	if keys.PrivKey != "" {
		keys.DecodePrivKey(a.cfg.GetApiSecret())
	}
	if keys.Passphrase != "" {
		keys.DecodePassKey(a.cfg.GetApiSecret())
	}

	return keys, nil
}
