package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/transport/http/handlers"
	"net/http"
	"time"
)

func NewHttpServer(cfg *config.Config, mux http.Handler) *http.Server {
	if mux == nil {
		panic("mux is nil")
	}

	return &http.Server{
		Addr:              cfg.GetApiPort(),
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}
}

func NewInternalMux() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/healthcheck"))

	return r
}

func NewExternalMux(apiKeyHandlers *handlers.ApiKeyHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/healthcheck"))

	r.Use(
		middleware.StripSlashes,
		middleware.RequestID,
	)

	r.Use(
		middleware.SetHeader("Access-Control-Allow-Origin", "*"),
		middleware.SetHeader("Access-Control-Allow-Methods", "POST"),
		middleware.SetHeader("Access-Control-Allow-Headers", "Content-Type"),
	)

	r.Post("/api-keys", apiKeyHandlers.AddApiKey)

	return r
}
