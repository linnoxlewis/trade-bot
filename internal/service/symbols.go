package service

import (
	"context"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

type SymbolRepo interface {
	AddSymbol(ctx context.Context, symbol string) error
	GetDefaultSymbols(ctx context.Context) (domain.SymbolList, error)
}

type SymbolService struct {
	orderRepo   OrderRepo
	symbolsRepo SymbolRepo
	logger      *log.Logger
}

func NewSymbolService(orderRepo OrderRepo, symbolsRepo SymbolRepo, logger *log.Logger) *SymbolService {
	return &SymbolService{
		orderRepo:   orderRepo,
		symbolsRepo: symbolsRepo,
		logger:      logger,
	}
}

func (s *SymbolService) GetActiveSymbols(ctx context.Context) (list domain.SymbolList, err error) {
	list, err = s.orderRepo.GetActiveSymbols(ctx)
	if err != nil {
		s.logger.ErrorLog.Println("err get active symbols: ", err)

		return nil, errors.InternalServerError(err)
	}

	return
}

func (s *SymbolService) GetDefaultSymbols(ctx context.Context) (list domain.SymbolList, err error) {
	list, err = s.symbolsRepo.GetDefaultSymbols(ctx)
	if err != nil {
		s.logger.ErrorLog.Println("err get default symbols: ", err)

		return nil, errors.InternalServerError(err)
	}

	return
}
