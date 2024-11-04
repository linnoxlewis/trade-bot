package service

import (
	"context"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

type SymbolRepo interface {
	AddSymbol(ctx context.Context, symbol string) error
}

type SymbolService struct {
	orderRepo OrderRepo
	logger    *log.Logger
}

func NewSymbolService(orderRepo OrderRepo, logger *log.Logger) *SymbolService {
	return &SymbolService{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

func (s *SymbolService) GetActiveSymbols(ctx context.Context) (list domain.SymbolList, err error) {
	list, err = s.orderRepo.GetActiveSymbols(ctx)
	if err != nil {
		s.logger.ErrorLog.Println("err get active symbols: ", err)

		return nil, err
	}

	return
}

func (s *SymbolService) GetDefaultSymbols(ctx context.Context) (list domain.SymbolList, err error) {
	return domain.SymbolList{
		"BTCUSDT",
		"ETHUSDT",
	}, nil
}
