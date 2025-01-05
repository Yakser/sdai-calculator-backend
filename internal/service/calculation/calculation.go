package service

import (
	"context"
	"errors"
	"sdai-calculator/internal/domain"
)

type CalculationService struct {
	repo domain.CalculationRepository
}

func NewCalculationService(repo domain.CalculationRepository) *CalculationService {
	return &CalculationService{repo: repo}
}

func (s *CalculationService) CalculateSDAI(
	ctx context.Context,
	painfulJoints int64,
	swollenJoints int64,
	patientActivityAssessment int64,
	physicianActivityAssessment int64,
	crp float64,
) (float64, error) {
	sdaiIndex := float64(painfulJoints) +
		float64(swollenJoints) +
		float64(patientActivityAssessment)/10 +
		float64(physicianActivityAssessment)/10 +
		crp

	_, err := s.repo.SaveCalculation(sdaiIndex)

	if err != nil {
		if errors.Is(err, domain.ErrCalculationAlreadyExists) {
			return 0, domain.ErrCalculationAlreadyExists
		}
		return 0, err
	}

	return sdaiIndex, nil
}

func (s *CalculationService) GetHistory(ctx context.Context, userID int64) ([]*domain.Calculation, error) {
	records, err := s.repo.GetCalculationsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return records, nil
}
