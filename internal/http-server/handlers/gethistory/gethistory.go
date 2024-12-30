package gethistory

import (
	"context"
	"log/slog"
	"sdai-calculator/internal/generated/server"
	"time"
)

type GetHistoryHandler struct {
	logger *slog.Logger
}

func NewGetHistoryHandler(logger *slog.Logger) *GetHistoryHandler {
	return &GetHistoryHandler{
		logger: logger,
	}
}

func (g *GetHistoryHandler) GetHistory(ctx context.Context, request server.GetHistoryRequestObject) (server.GetHistoryResponseObject, error) {
	// todo
	return server.GetHistory200JSONResponse{
		GetHistoryResponseJSONResponse: server.GetHistoryResponseJSONResponse{
			History: []server.SDAIRecord{
				{
					MeasureDatetime: time.Now(),
					SdaiIndex:       1,
				},
			},
		},
	}, nil
}
