package tests

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	v1 "github.com/linnoxlewis/trade-bot/internal/transport/api/controller/v1"
	"github.com/linnoxlewis/trade-bot/internal/transport/api/controller/v1/tests/mocks"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestApiKeyController_AddApiKey(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
	}{
		{
			name: "ValidRequest",
			requestBody: dto.ApiKeys{
				PubKey:     "test",
				PrivKey:    "test",
				PassPhrase: "TestTest",
				Exchange:   "Binance",
				UserId:     123,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "InvalidRequestPub",
			requestBody: dto.ApiKeys{
				PrivKey:    "test",
				PassPhrase: "TestTest",
				Exchange:   "Binance",
				UserId:     123,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "InvalidRequestPriv",
			requestBody: dto.ApiKeys{
				PubKey:     "",
				PassPhrase: "TestTest",
				Exchange:   "Binance",
				UserId:     123,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "InvalidRequestExchange",
			requestBody: dto.ApiKeys{
				PubKey:     "",
				PrivKey:    "test",
				PassPhrase: "TestTest",
				UserId:     123,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "InvalidRequestUserId",
			requestBody: dto.ApiKeys{
				PubKey:     "",
				PrivKey:    "test",
				PassPhrase: "TestTest",
				Exchange:   "Binance",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockApiKeyService(ctrl)
	mockLogger := log.NewLogger()
	apiKeyController := v1.NewApiKeyController(mockService, mockLogger)

	router.POST("api/v1/api-key/", apiKeyController.AddApiKey)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(http.MethodPost, "api/v1/api-key/", bytes.NewBuffer(requestBody))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
