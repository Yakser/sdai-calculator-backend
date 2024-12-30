package handlers

import (
	"log/slog"
	"sdai-calculator/internal/generated/server"
	"sdai-calculator/internal/http-server/handlers/calculate"
	"sdai-calculator/internal/http-server/handlers/gethistory"

	//"sdai-calculator/internal/http-server/handlers/redirect"
	//"sdai-calculator/internal/http-server/handlers/url/save"
	"sdai-calculator/internal/storage"
)

var _ server.StrictServerInterface = (*Handlers)(nil)

type Handlers struct {
	*calculate.CalculateHandler
	*gethistory.GetHistoryHandler
}

func NewHandlers(logger *slog.Logger, store storage.Storage) *Handlers {
	return &Handlers{
		CalculateHandler:  calculate.NewCalculateHandler(logger),
		GetHistoryHandler: gethistory.NewGetHistoryHandler(logger),
	}
}
