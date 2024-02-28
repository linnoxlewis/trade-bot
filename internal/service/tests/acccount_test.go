package tests

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	srvErr "github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/internal/service"
	mock_service "github.com/linnoxlewis/trade-bot/internal/service/tests/mocks"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountService_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiKeyRepo := mock_service.NewMockApiKeyRepo(ctrl)
	mockExchanger := mock_service.NewMockExchanger(ctrl)
	mockLogger := log.NewLogger()
	cfg := &config.Config{}
	i18nSrv := i18n.NewI18n("../../../"+consts.LocaleDataPath, helper.GetLanguageList())
	accountService := service.NewAccountService(cfg, mockApiKeyRepo, mockExchanger, i18nSrv, mockLogger)

	testCases := []struct {
		name              string
		apiKeys           *domain.ApiKeys
		exchange          string
		expectedBalance   domain.Balance
		expectedError     error
		expectedErrorCode int
		apiKeyRepoError   error
		exchangerError    error
		apiKeysRepoError  error
		emptyApiKeys      bool
		invalidPrivKey    bool
		invalidPassKey    bool
	}{
		{
			name:              "Get balance success",
			apiKeys:           &domain.ApiKeys{},
			exchange:          "exchangeName",
			expectedBalance:   domain.Balance{},
			expectedError:     nil,
			expectedErrorCode: 200,
			apiKeyRepoError:   nil,
			exchangerError:    nil,
			apiKeysRepoError:  nil,
			emptyApiKeys:      false,
			invalidPrivKey:    false,
			invalidPassKey:    false,
		},
		{
			name:              "Error getting API keys",
			apiKeys:           nil,
			exchange:          "exchangeName",
			expectedBalance:   nil,
			expectedError:     srvErr.InternalServerError(errors.New("error getting API keys")),
			expectedErrorCode: 500,
			apiKeyRepoError:   errors.New("error getting API keys"),
			exchangerError:    nil,
			apiKeysRepoError:  nil,
			emptyApiKeys:      true,
			invalidPrivKey:    false,
			invalidPassKey:    false,
		},
		{
			name:              "Empty exchange name",
			apiKeys:           nil,
			exchange:          "",
			expectedBalance:   nil,
			expectedError:     srvErr.BadRequestError("empty exchange name"),
			expectedErrorCode: 400,
			apiKeyRepoError:   nil,
			exchangerError:    nil,
			apiKeysRepoError:  nil,
			emptyApiKeys:      true,
			invalidPrivKey:    false,
			invalidPassKey:    false,
		},
		{
			name:              "Invalid private key",
			apiKeys:           &domain.ApiKeys{},
			exchange:          "exchangeName",
			expectedBalance:   nil,
			expectedError:     srvErr.BadRequestError("invalid private key"),
			expectedErrorCode: 400,
			apiKeyRepoError:   nil,
			exchangerError:    nil,
			apiKeysRepoError:  nil,
			emptyApiKeys:      false,
			invalidPrivKey:    true,
			invalidPassKey:    false,
		},
		{
			name:              "Invalid passphrase key",
			apiKeys:           &domain.ApiKeys{},
			exchange:          "exchangeName",
			expectedBalance:   nil,
			expectedError:     srvErr.BadRequestError("invalid passphrase key"),
			expectedErrorCode: 400,
			apiKeyRepoError:   nil,
			exchangerError:    nil,
			apiKeysRepoError:  nil,
			emptyApiKeys:      false,
			invalidPrivKey:    false,
			invalidPassKey:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockApiKeyRepo.EXPECT().GetApiKeysByUserIdAndExchange(gomock.Any(), gomock.Any(), tc.exchange).Return(tc.apiKeys, tc.apiKeysRepoError)
			if !tc.emptyApiKeys {
				if tc.invalidPrivKey {
					tc.apiKeys.DecodePrivKey("")
				}
				if tc.invalidPassKey {
					tc.apiKeys.DecodePassKey("")
				}
			}
			if tc.apiKeys != nil && !tc.emptyApiKeys {
				mockExchanger.EXPECT().Balance(tc.apiKeys, tc.exchange).Return(tc.expectedBalance, tc.exchangerError)
			}

			balance, err := accountService.GetBalance(context.Background(), 123, tc.exchange)
			if err != nil {
				resultErr := err.(srvErr.Error)
				assert.Equal(t, tc.expectedError, resultErr)
				assert.Equal(t, tc.expectedErrorCode, resultErr.GetCode())
			} else {
				assert.Equal(t, tc.expectedBalance, balance)
			}
		})
	}
}
