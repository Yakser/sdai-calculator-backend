package domain

type CalculationRepository interface {
	SaveCalculation(sdaiIndex float64) (int64, error)
	GetCalculationsByUserID(userID int64) ([]*Calculation, error)
}
