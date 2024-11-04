package handlers

import (
	"context"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/internal/pkg/validation"
	"github.com/linnoxlewis/trade-bot/internal/transport/http/helper"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"net/http"
)

type ApiKeyService interface {
	AddApiKeys(ctx context.Context, apiKeys *dto.ApiKeys) error
}

type ApiKeyHandler struct {
	apiKeySrv ApiKeyService
	logger    *log.Logger
	validator validation.Validator
}

func NewApiKeyHandler(apiKeySrv ApiKeyService, logger *log.Logger) *ApiKeyHandler {
	return &ApiKeyHandler{
		apiKeySrv: apiKeySrv,
		logger:    logger,
	}
}

// AddApiKey api добавления апи ключей
func (a *ApiKeyHandler) AddApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request dto.ApiKeys
	if err := helper.ParseJSON(r.Body, &request); err != nil {
		helper.SendBadRequestResponse(w, errors.BadRequestError("failed to parse request body"))

		return
	}

	if err := a.validator.ValidateAddApiKeys(&request); err != nil {
		helper.SendBadRequestResponse(w, errors.BadRequestError(err.Error()))

		return
	}

	if err := a.apiKeySrv.AddApiKeys(ctx, &request); err != nil {
		helper.SendBadRequestResponse(w, errors.BadRequestError(err.Error()))

		return
	}

	helper.SendSuccessResponse(w, "API key(s) added successfully")
}

func (a *ApiKeyHandler) RemoveApiKey(w http.ResponseWriter, r http.Request) {
	helper.SendSuccessResponse(w, "API key(s) removed")
}
