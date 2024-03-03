package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

type ApiKeyService interface {
	AddApiKeys(ctx context.Context, apiKeys *dto.ApiKeys) error
	DeleteApiKey(ctx context.Context, userId int64, exchange string) error
	ClearApiKey(ctx context.Context, userId int64, exchange string) error
	GetApiKeyByExchangeAndId(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error)
}

type ApiKeyController struct {
	apiKeySrv ApiKeyService
	logger    *log.Logger
}

func NewApiKeyController(apiKeySrv ApiKeyService, logger *log.Logger) *ApiKeyController {
	return &ApiKeyController{
		apiKeySrv: apiKeySrv,
		logger:    logger,
	}
}

func (a *ApiKeyController) AddApiKey(c *gin.Context) {
	rqt := &dto.ApiKeys{}
	if err := c.BindJSON(rqt); err != nil {
		helper.JsonErrorResponse(c)

		return
	}

	if err := rqt.Validate(); err != nil {
		helper.BadRequestErrorResponse(c, err)

		return
	}

	if err := a.apiKeySrv.AddApiKeys(c, rqt); err != nil {
		helper.BadRequestErrorResponse(c, err)

		return
	}

	helper.SuccessResponse(c, "apiKeys added")
}

func (a *ApiKeyController) RemoveApiKey(c *gin.Context) {
	helper.SuccessResponse(c, "apiKeys removed")
}
