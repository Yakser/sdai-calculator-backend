package calculate

import (
	"context"
	"log/slog"
	"sdai-calculator/internal/generated/server"
)

type CalculateHandler struct {
	logger *slog.Logger
}

func NewCalculateHandler(logger *slog.Logger) *CalculateHandler {
	return &CalculateHandler{
		logger: logger,
	}
}

func (c *CalculateHandler) Calculate(ctx context.Context, request server.CalculateRequestObject) (server.CalculateResponseObject, error) {
	// todo
	return server.Calculate200JSONResponse{
		CalculateResponseJSONResponse: server.CalculateResponseJSONResponse{SdaiIndex: 1},
	}, nil
}
