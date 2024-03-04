package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	v1 "github.com/linnoxlewis/trade-bot/internal/transport/api/controller/v1"
)

type ApiKeyService interface {
	AddApiKeys(ctx context.Context, apiKeys *dto.ApiKeys) error
	DeleteApiKey(ctx context.Context, userId int64, exchange string) error
}

func RegisterApiKeyRoutes(engine *gin.Engine, orderCtrl *v1.ApiKeyController) {
	group := engine.Group("api/v1/api-key")
	group.POST("/", orderCtrl.AddApiKey)
	group.DELETE("/{exchange}}", orderCtrl.RemoveApiKey)
}
