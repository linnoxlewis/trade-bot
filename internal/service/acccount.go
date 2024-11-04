package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/pkg/exchanger"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

var (
	ErrApiKeysNotFound    = errors.New("api keys not found")
	ErrBalanceFetchFailed = errors.New("failed to fetch balance from exchange")
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
		logger:     logger,
	}
}

// GetBalance получает баланс пользователя на бирже
func (a *AccountService) GetBalance(ctx context.Context, userId int64, exchange string) (domain.Balances, error) {
	keys, err := a.getApiKeys(ctx, userId, exchange)
	if err != nil {
		return nil, err
	}

	result, err := a.exchanger.Balance(keys, exchange)
	if err != nil {
		a.logger.ErrorLog.Printf("failed to fetch balance for user %d on exchange %s: %v", userId, exchange, err)

		return nil, fmt.Errorf("%w: %v", ErrBalanceFetchFailed, err)
	}

	return result, nil
}

// getApiKeys получает и расшифровывает API-ключи для пользователя
func (a *AccountService) getApiKeys(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error) {
	keys, err := a.apiKeyRepo.GetApiKeysByUserIdAndExchange(ctx, userId, exchange)
	if err != nil {
		a.logger.ErrorLog.Printf("failed to retrieve API keys for user %d on exchange %s: %v", userId, exchange, err)

		return nil, fmt.Errorf("repository error: %w", err)
	}

	if keys == nil {
		a.logger.ErrorLog.Printf("no API keys found for user %d on exchange %s", userId, exchange)

		return nil, fmt.Errorf("%w: userId %d, exchange %s", ErrApiKeysNotFound, userId, exchange)
	}

	if err := a.decryptApiKeys(keys); err != nil {
		a.logger.ErrorLog.Printf("failed to decrypt API keys for user %d on exchange %s: %v", userId, exchange, err)

		return nil, fmt.Errorf("decryption error: %w", err)
	}

	return keys, nil
}

// decryptApiKeys расшифровывает приватный ключ и passphrase, если они присутствуют
func (a *AccountService) decryptApiKeys(keys *domain.ApiKeys) error {
	if keys.PrivKey != "" {
		if err := keys.DecryptPrivKey(a.cfg.GetApiSecret()); err != nil {
			return fmt.Errorf("failed to decrypt private key: %w", err)
		}
	}

	if keys.Passphrase != "" {
		if err := keys.DecryptPassphrase(a.cfg.GetApiSecret()); err != nil {
			return fmt.Errorf("failed to decrypt passphrase: %w", err)
		}
	}

	return nil
}
