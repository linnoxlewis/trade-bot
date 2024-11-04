package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEncryptionError = errors.New("encryption error")
)

type ApiKeyRepo interface {
	AddApiKeys(ctx context.Context, keys *domain.ApiKeys) error
	DeleteApiKey(ctx context.Context, userId int64, exchange string) error
	ClearApiKey(ctx context.Context, userId int64) error
	GetApiKeysByUserIdAndExchange(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error)
	GetApiKeysByUserId(ctx context.Context, userId int64) ([]*domain.ApiKeys, error)
}

type ApiKeysService struct {
	cfg        *config.Config
	apiKeyRepo ApiKeyRepo
	userRepo   UserRepo
	i18n       *i18n.I18n
	logger     *log.Logger
}

func NewApiKeysService(cfg *config.Config, apiKeyRepo ApiKeyRepo, userRepo UserRepo, i18n *i18n.I18n, logger *log.Logger) *ApiKeysService {
	return &ApiKeysService{
		cfg:        cfg,
		apiKeyRepo: apiKeyRepo,
		userRepo:   userRepo,
		i18n:       i18n,
		logger:     logger,
	}
}

// AddApiKeys adds new API keys for a user after decrypting and encrypting them appropriately.
func (a *ApiKeysService) AddApiKeys(ctx context.Context, apiKeys *dto.ApiKeys) error {
	if !a.userRepo.ExistUser(ctx, apiKeys.UserId) {
		return fmt.Errorf("%w: user ID %d", ErrUserNotFound, apiKeys.UserId)
	}

	encryptedPrivateKey, passphrase, err := a.processKeys(apiKeys)
	if err != nil {
		return err
	}

	userApiKey := domain.NewApiKeys(apiKeys.UserId, apiKeys.Exchange, apiKeys.PubKey, encryptedPrivateKey, passphrase)

	if err := a.apiKeyRepo.AddApiKeys(ctx, userApiKey); err != nil {
		a.logError("failed to add API keys to repository", err)
		return fmt.Errorf("repository error: %w", err)
	}

	return nil
}

// DeleteApiKey removes API keys for a user by exchange.
func (a *ApiKeysService) DeleteApiKey(ctx context.Context, userId int64, exchange string) error {
	if err := a.apiKeyRepo.DeleteApiKey(ctx, userId, exchange); err != nil {
		a.logError(fmt.Sprintf("failed to delete API key for user %d and exchange %s", userId, exchange), err)

		return fmt.Errorf("repository error: %w", err)
	}
	return nil
}

// ClearApiKey removes all API keys for a user.
func (a *ApiKeysService) ClearApiKey(ctx context.Context, userId int64) error {
	if err := a.apiKeyRepo.ClearApiKey(ctx, userId); err != nil {
		a.logError(fmt.Sprintf("failed to clear API keys for user %d", userId), err)

		return fmt.Errorf("repository error: %w", err)
	}
	return nil
}

// GetApiKeyByExchangeAndId retrieves API keys for a user by exchange.
func (a *ApiKeysService) GetApiKeyByExchangeAndId(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error) {
	apiKeys, err := a.apiKeyRepo.GetApiKeysByUserIdAndExchange(ctx, userId, exchange)
	if err != nil {
		a.logError(fmt.Sprintf("failed to retrieve API key for user %d and exchange %s", userId, exchange), err)

		return nil, fmt.Errorf("repository error: %w", err)
	}
	return apiKeys, nil
}

// processKeys handles decryption and encryption of API keys.
func (a *ApiKeysService) processKeys(apiKeys *dto.ApiKeys) (string, string, error) {
	privateKey, err := helper.DecryptClientMessage(apiKeys.PrivKey, a.cfg.GetClientSecretKey())
	if err != nil {
		a.logError("failed to decrypt private key", err)

		return "", "", fmt.Errorf("%w: failed to decrypt private key", ErrEncryptionError)
	}

	encryptedPrivateKey, err := helper.EncryptMessage(privateKey, a.cfg.GetApiSecret())
	if err != nil {
		a.logError("failed to encrypt private key", err)

		return "", "", fmt.Errorf("%w: failed to encrypt private key", ErrEncryptionError)
	}

	var encryptedPassphrase string
	if !apiKeys.EmptyPassPhrase() {
		passphrase, err := helper.DecryptClientMessage(apiKeys.PassPhrase, a.cfg.GetClientSecretKey())
		if err != nil {
			a.logError("failed to decrypt passphrase", err)

			return "", "", fmt.Errorf("%w: failed to decrypt passphrase", ErrEncryptionError)
		}

		encryptedPassphrase, err = helper.EncryptMessage(passphrase, a.cfg.GetApiSecret())
		if err != nil {
			a.logError("failed to encrypt passphrase", err)

			return "", "", fmt.Errorf("%w: failed to encrypt passphrase", ErrEncryptionError)
		}
	}

	return encryptedPrivateKey, encryptedPassphrase, nil
}

// logError is a helper method for logging errors with additional context.
func (a *ApiKeysService) logError(message string, err error) {
	a.logger.ErrorLog.Printf("%s: %v", message, err)
}
