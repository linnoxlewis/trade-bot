package server

import (
	"context"
	"github.com/gin-gonic/gin"
	ctrl "github.com/linnoxlewis/trade-bot/internal/transport/api/controller/v1"
	v1 "github.com/linnoxlewis/trade-bot/internal/transport/api/route/v1"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

type TransportInterface interface {
	StartServer()
	StopServer()
}

type ApiServer struct {
	server *http.Server
	logger *log.Logger
}

func NewApiServer(
	port string,
	serverMode string,
	apikeySrv ctrl.ApiKeyService,
	logger *log.Logger,
) *ApiServer {
	engine := gin.Default()
	gin.SetMode(serverMode)
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
	engine.GET("/server-time", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": time.Now().Unix(),
		})
	})
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiKeyCtrl := ctrl.NewApiKeyController(apikeySrv, logger)
	v1.RegisterApiKeyRoutes(engine, apiKeyCtrl)
	server := &http.Server{
		Addr:     port,
		Handler:  engine,
		ErrorLog: logger.ErrorLog,
	}

	return &ApiServer{server: server, logger: logger}
}

func (r *ApiServer) StartServer() {
	r.logger.InfoLog.Println("API Server starting...")

	if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		r.logger.ErrorLog.Fatalf("Failed to listen and serve: %+v", err)
	}
}

func (r *ApiServer) StopServer() {
	r.logger.InfoLog.Println("API Server stopping...")

	if err := r.server.Shutdown(context.Background()); err != nil {
		r.logger.ErrorLog.Fatalf("Failed stopped serve: %+v", err)
	}
}
