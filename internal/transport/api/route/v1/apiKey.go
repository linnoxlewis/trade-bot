package v1

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/linnoxlewis/trade-bot/internal/transport/api/controller/v1"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

func RegisterApiKeyRoutes(
	engine *gin.Engine,
	orderSrv v1.ApiKeyService,
	logger *log.Logger,
) {
	group := engine.Group("api/v1/api-key")
	orderCtrl := v1.NewApiKeyController(orderSrv, logger)

	group.POST("/", orderCtrl.AddApiKey)
	group.DELETE("/{exchange}}", orderCtrl.RemoveApiKey)
}
