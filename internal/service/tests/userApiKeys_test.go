package tests

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/internal/service"
	mock_service "github.com/linnoxlewis/trade-bot/internal/service/tests/mocks"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiKeysService_AddApiKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiKeyRepo := mock_service.NewMockApiKeyRepo(ctrl)
	mockUserRepo := mock_service.NewMockUserRepo(ctrl)
	mockLogger := log.NewLogger()
	cfg := &config.Config{}
	i18nSrv := i18n.NewI18n(consts.LocaleDataPath, helper.GetLanguageList())
	apiKeyService := service.NewApiKeysService(cfg, mockApiKeyRepo, mockUserRepo, i18nSrv, mockLogger)

	testCases := []struct {
		name              string
		apiKeys           *dto.ApiKeys
		expectedError     error
		expectedErrorCode int
		apiKeyRepoError   error
	}{
		{
			name:              "AddApiKeys success",
			apiKeys:           &dto.ApiKeys{UserId: 1, Exchange: "TestExchange", PubKey: "TestPubKey", PrivKey: "TestPrivKey", PassPhrase: "TestPassPhrase"},
			expectedError:     nil,
			expectedErrorCode: 0,
			apiKeyRepoError:   nil,
		},
		{
			name:              "AddApiKeys missing user",
			apiKeys:           &dto.ApiKeys{UserId: 0, Exchange: "TestExchange", PubKey: "TestPubKey", PrivKey: "TestPrivKey", PassPhrase: "TestPassPhrase"},
			expectedError:     errors.BadRequestError("userNotFound"),
			expectedErrorCode: 400,
			apiKeyRepoError:   nil,
		},
		// Добавьте дополнительные сценарии тестирования
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo.EXPECT().ExistUser(gomock.Any(), gomock.Any()).Return(tc.apiKeys.UserId != 0)
			if tc.apiKeys.UserId != 0 {
				mockApiKeyRepo.EXPECT().AddApiKeys(gomock.Any(), gomock.Any()).Return(tc.apiKeyRepoError)
			}

			err := apiKeyService.AddApiKeys(context.Background(), tc.apiKeys)
			if err != nil {
				resultErr := err.(errors.Error)
				assert.Equal(t, tc.expectedError, resultErr)
				assert.Equal(t, tc.expectedErrorCode, resultErr.GetCode())
			}
		})
	}
}
