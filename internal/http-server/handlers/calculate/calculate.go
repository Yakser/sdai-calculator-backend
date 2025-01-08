package calculate

import (
	"context"
	"log/slog"
	"sdai-calculator/internal/domain"
	"sdai-calculator/internal/generated/server"
	"sdai-calculator/internal/lib/logger/sl"
)

type CalculateHandler struct {
	logger     *slog.Logger
	calculator Calculator
}

type Calculator interface {
	CalculateSDAI(ctx context.Context, painfulJoints int64, swollenJoints int64, patientActivityAssessment int64, physicianActivityAssessment int64, crp float64) (float64, error)
}

func NewCalculateHandler(logger *slog.Logger, calculator Calculator) *CalculateHandler {
	return &CalculateHandler{
		logger:     logger,
		calculator: calculator,
	}
}

func (c *CalculateHandler) Calculate(ctx context.Context, request server.CalculateRequestObject) (server.CalculateResponseObject, error) {
	params := request.Body.Parameters

	sdaiIndex, err := c.calculator.CalculateSDAI(ctx,
		params.PainfulJoints,
		params.SwollenJoints,
		params.PatientActivityAssessment,
		params.PhysicianActivityAssessment,
		params.Crp,
	)
	if err != nil {
		if err == domain.ErrCalculationAlreadyExists {
			return server.Calculate400JSONResponse{
				CalculateErrorResponseJSONResponse: server.CalculateErrorResponseJSONResponse{
					Message: "Calculation already exists",
				},
			}, nil
		}

		c.logger.Error("failed to calculate SDAI index", sl.Err(err))

		// fixme: change to 500 err
		return server.Calculate400JSONResponse{
			CalculateErrorResponseJSONResponse: server.CalculateErrorResponseJSONResponse{
				Message: "Internal server error",
			},
		}, nil
	}

	return server.Calculate200JSONResponse{
		CalculateResponseJSONResponse: server.CalculateResponseJSONResponse{
			SdaiIndex: sdaiIndex,
		},
	}, nil
}
