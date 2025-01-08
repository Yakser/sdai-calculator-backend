package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"sdai-calculator/internal/domain"
)

type CalculationRepository struct {
	db *sqlx.DB
}

func NewCalculationRepository(connectionString string) (*CalculationRepository, error) {
	const component = "storage.postgresql.NewCalculationStorage"

	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", component, err)
	}

	// Проверяем соединение с базой данных.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%v: %w", component, err)
	}

	return &CalculationRepository{
		db: db,
	}, nil
}

func (s *CalculationRepository) SaveCalculation(
	painfulJoints int64,
	swollenJoints int64,
	patientActivityAssessment int64,
	physicianActivityAssessment int64,
	crp float64,
	sdaiIndex float64,
) (int64, error) {
	const component = "storage.postgresql.SaveCalculation"

	// todo: use sqlx
	query, err := s.db.Prepare(`
		INSERT INTO calculations (user_id, painful_joints, swollen_joints, physician_activity_assessment, patient_activity_assessment, crp, sdai_index) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`)
	if err != nil {
		return 0, fmt.Errorf("%v: %w", component, err)
	}
	defer query.Close()

	// fixme
	userID := -1

	var id int64
	err = query.QueryRow(userID,
		painfulJoints,
		swollenJoints,
		patientActivityAssessment,
		physicianActivityAssessment,
		crp, fmt.Sprintf("%v", sdaiIndex)).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// fixme: оно тут не появится
			if pgErr.Code == pgerrcode.UniqueViolation {
				return 0, fmt.Errorf("%v: %w", component, domain.ErrCalculationAlreadyExists)
			}
		}
		return 0, fmt.Errorf("%v: %w", component, err)
	}

	return id, nil
}

func (s *CalculationRepository) GetCalculationsByUserID(ctx context.Context, userID int64) ([]domain.Calculation, error) {
	const component = "storage.postgresql.GetCalculationsByUserID"
	// todo: pagination
	// todo
	const getCalculationsQuery = `
		SELECT id, user_id, painful_joints, swollen_joints, physician_activity_assessment, patient_activity_assessment, crp, sdai_index, created_at
		FROM calculations
		WHERE user_id = $1
	`

	query, args, err := sqlx.In(getCalculationsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", component, err)
	}

	query = s.db.Rebind(query)

	var calc []domain.Calculation

	err = s.db.SelectContext(ctx, &calc, query, args...)

	return calc, err
}
