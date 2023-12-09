package tests_test

import (
	"context"
	"errors"
	srvErr "github.com/linnoxlewis/trade-bot/internal/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/service"
	mock_service "github.com/linnoxlewis/trade-bot/internal/service/tests/repo_mocks"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestSymbolService_GetActiveSymbols(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock_service.NewMockOrderRepo(ctrl)
	mockLogger := log.NewLogger()

	symbolService := service.NewSymbolService(mockOrderRepo, nil, mockLogger)

	testCases := []struct {
		name             string
		repoResponse     domain.SymbolList
		repoError        error
		expectedResponse domain.SymbolList
		expectedError    error
	}{
		{
			name:             "Success case",
			repoResponse:     domain.SymbolList{"BTC_USDT", "BNB_USDT"},
			repoError:        nil,
			expectedResponse: domain.SymbolList{"BTC_USDT", "BNB_USDT"},
			expectedError:    nil,
		},
		{
			name:             "Error case",
			repoResponse:     nil,
			repoError:        errors.New("order repository error"),
			expectedResponse: nil,
			expectedError:    srvErr.InternalServerError(errors.New("order repository error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockOrderRepo.EXPECT().GetActiveSymbols(gomock.Any()).Return(tc.repoResponse, tc.repoError)

			response, err := symbolService.GetActiveSymbols(context.Background())

			assert.Equal(t, tc.expectedResponse, response)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestSymbolService_GetDefaultSymbols(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSymbolRepo := mock_service.NewMockSymbolRepo(ctrl)
	mockLogger := log.NewLogger()

	symbolService := service.NewSymbolService(nil, mockSymbolRepo, mockLogger)

	testCases := []struct {
		name             string
		repoResponse     domain.SymbolList
		repoError        error
		expectedResponse domain.SymbolList
		expectedError    error
	}{
		{
			name:             "Success case",
			repoResponse:     domain.SymbolList{"BTC_USDT", "BNB_USDT"},
			repoError:        nil,
			expectedResponse: domain.SymbolList{"BTC_USDT", "BNB_USDT"},
			expectedError:    nil,
		},
		{
			name:             "Error case",
			repoResponse:     nil,
			repoError:        errors.New("symbol repository error"),
			expectedResponse: nil,
			expectedError:    srvErr.InternalServerError(errors.New("symbol repository error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSymbolRepo.EXPECT().GetDefaultSymbols(gomock.Any()).Return(tc.repoResponse, tc.repoError)

			response, err := symbolService.GetDefaultSymbols(context.Background())

			assert.Equal(t, tc.expectedResponse, response)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
