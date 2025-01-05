package handlers

import (
	"context"
	"log/slog"
	"sdai-calculator/internal/domain"
	"sdai-calculator/internal/generated/server"
	"sdai-calculator/internal/http-server/handlers/calculate"
	"sdai-calculator/internal/http-server/handlers/gethistory"
)

var _ server.StrictServerInterface = (*Handlers)(nil)

// fixme: убрать из имен хендлеров имя их пакетов

type Handlers struct {
	*calculate.CalculateHandler
	*gethistory.GetHistoryHandler
}

type CalculationService interface {
	CalculateSDAI(ctx context.Context, painfulJoints int64, swollenJoints int64, patientActivityAssessment int64, physicianActivityAssessment int64, crp float64) (float64, error)
	GetHistory(ctx context.Context, userID int64) ([]*domain.Calculation, error)
}

func NewHandlers(logger *slog.Logger, calculationService CalculationService) *Handlers {
	return &Handlers{
		CalculateHandler:  calculate.NewCalculateHandler(logger, calculationService),
		GetHistoryHandler: gethistory.NewGetHistoryHandler(logger, calculationService),
	}
}
