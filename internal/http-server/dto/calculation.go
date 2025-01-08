package dto

import (
	"sdai-calculator/internal/domain"
	"sdai-calculator/internal/generated/server"
	"strconv"
)

func ToHistoryRecord(calculation domain.Calculation) server.HistoryRecord {
	floatSdaiIndex, err := strconv.ParseFloat(calculation.SdaiIndex, 64)
	// should never happen
	if err != nil {
		floatSdaiIndex = -1
	}

	return server.HistoryRecord{
		SdaiIndex:       floatSdaiIndex,
		MeasureDatetime: calculation.CreatedAt,
		Parameters: server.CalculationParameters{
			Crp:                         calculation.Crp,
			PainfulJoints:               calculation.PainfulJoints,
			PatientActivityAssessment:   calculation.PatientActivityAssessment,
			PhysicianActivityAssessment: calculation.PhysicalActivityAssessment,
			SwollenJoints:               calculation.SwollenJoints,
		},
	}
}

func ToHistoryRecords(calculations []domain.Calculation) []server.HistoryRecord {
	result := make([]server.HistoryRecord, len(calculations))
	for i, c := range calculations {
		result[i] = ToHistoryRecord(c)
	}
	return result
}
