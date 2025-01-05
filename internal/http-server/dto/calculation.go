package dto

import (
	"sdai-calculator/internal/domain"
	"sdai-calculator/internal/generated/server"
	"strconv"
)

func ToSDAIRecord(calculation *domain.Calculation) server.SDAIRecord {
	floatSdaiIndex, err := strconv.ParseFloat(calculation.SdaiIndex, 64)
	// should never happen
	if err != nil {
		floatSdaiIndex = -1
	}

	return server.SDAIRecord{
		SdaiIndex:       floatSdaiIndex,
		MeasureDatetime: calculation.CreatedAt,
	}
}

func ToSDAIRecords(calculations []*domain.Calculation) []server.SDAIRecord {
	result := make([]server.SDAIRecord, len(calculations))
	for i, c := range calculations {
		result[i] = ToSDAIRecord(c)
	}
	return result
}
