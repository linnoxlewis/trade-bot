package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

type ApiKeyService interface {
	AddApiKeys(ctx context.Context, userId uuid.UUID, apiKeys *dto.ApiKeys) error
	DeleteApiKey(ctx context.Context, userId uuid.UUID, exchange string) error
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

// SetApiKey  godoc
// @Summary         Управление API ключами
// @Description     Добавляет binance api ключ
// @Tags            BinanceProvider
// @Param  			pubKey  query  string  true  "Публичный ключь"
// @Param  			privKey  query  integer  true   "Приватный ключь"
// @Param  			description  query  integer  true   "Имя ключа"
// @Accept          json
// @Produce         json
// @Success      200  {object}  form.ApiKey "Ключ добавлен"
// @Failure      500  {object}  response.InternalServerError      "Произошла внутренняя ошибка сервера"
// @Router       /api/v1/binance-provider/set-api-key [post]
/*
func (o *OrderController) CreateFiatPayment(c *gin.Context) {
	rqt := &dto.FiatPayment{}
	if err := c.BindJSON(rqt); err != nil {
		helper.JsonErrorResponse(c)

		return
	}

	if err := rqt.Validate(); err != nil {
		helper.BadRequestErrorResponse(c, err)
		fmt.Println(err)
		return
	}

	if err := o.orderSrv.CreateFiatPayment(c, rqt); err != nil {
		helper.BadRequestErrorResponse(c, err)

		fmt.Println(err)
		return
	}
	helper.SuccessResponse(c, "payment completed")
}*/

func (a *ApiKeyController) AddApiKey(c *gin.Context) {
	helper.SuccessResponse(c, "apiKeys added")
}

func (a *ApiKeyController) RemoveApiKey(c *gin.Context) {
	helper.SuccessResponse(c, "apiKeys removed")
}
