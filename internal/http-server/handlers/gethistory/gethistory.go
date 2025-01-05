package gethistory

import (
	"context"
	"log/slog"
	"sdai-calculator/internal/domain"
	"sdai-calculator/internal/generated/server"
	"sdai-calculator/internal/http-server/dto"
	"sdai-calculator/internal/lib/logger/sl"
)

type GetHistoryHandler struct {
	logger        *slog.Logger
	historyGetter HistoryGetter
}

type HistoryGetter interface {
	GetHistory(ctx context.Context, userID int64) ([]*domain.Calculation, error)
}

func NewGetHistoryHandler(logger *slog.Logger, historyGetter HistoryGetter) *GetHistoryHandler {
	return &GetHistoryHandler{
		logger:        logger,
		historyGetter: historyGetter,
	}
}

func (g *GetHistoryHandler) GetHistory(ctx context.Context, request server.GetHistoryRequestObject) (server.GetHistoryResponseObject, error) {
	// fixme
	userID := int64(-1)
	calculations, err := g.historyGetter.GetHistory(ctx, userID)

	if err != nil {
		g.logger.Error("failed to get history", sl.Err(err))

		return server.GetHistory400JSONResponse{
			GetHistoryErrorResponseJSONResponse: server.GetHistoryErrorResponseJSONResponse{
				Code:    nil,
				Message: "failed to get history",
			},
		}, nil
	}

	return server.GetHistory200JSONResponse{
		GetHistoryResponseJSONResponse: server.GetHistoryResponseJSONResponse{
			History: dto.ToSDAIRecords(calculations),
		},
	}, nil
}
