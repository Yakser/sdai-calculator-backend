package domain

import "context"

type CalculationRepository interface {
	SaveCalculation(
		painfulJoints int64,
		swollenJoints int64,
		patientActivityAssessment int64,
		physicianActivityAssessment int64,
		crp float64,
		sdaiIndex float64,
	) (int64, error)
	GetCalculationsByUserID(ctx context.Context, userID int64) ([]Calculation, error)
}
