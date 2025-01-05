package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"sdai-calculator/internal/domain"
)

type CalculationRepository struct {
	db *sql.DB
}

func NewCalculationRepository(connectionString string) (*CalculationRepository, error) {
	const component = "storage.postgresql.NewCalculationStorage"

	db, err := sql.Open("pgx", connectionString)
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

func (s *CalculationRepository) SaveCalculation(sdaiIndex float64) (int64, error) {
	const component = "storage.postgresql.SaveCalculation"

	query, err := s.db.Prepare(`
		INSERT INTO calculations (user_id, sdai_index) 
		VALUES ($1, $2) RETURNING id
	`)
	if err != nil {
		return 0, fmt.Errorf("%v: %w", component, err)
	}
	defer query.Close()

	var id int64
	// fixme: change user id
	err = query.QueryRow(-1, fmt.Sprintf("%v", sdaiIndex)).Scan(&id)
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

func (s *CalculationRepository) GetCalculationsByUserID(userID int64) ([]*domain.Calculation, error) {
	const component = "storage.postgresql.GetCalculationsByUserID"

	query, err := s.db.Prepare(`
		SELECT id, user_id, sdai_index, created_at
		FROM calculations
		WHERE user_id = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", component, err)
	}
	defer query.Close()

	rows, err := query.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", component, err)
	}
	defer rows.Close()

	var calculations []*domain.Calculation
	for rows.Next() {
		var calculation domain.Calculation
		if err := rows.Scan(&calculation.ID, &calculation.UserID, &calculation.SdaiIndex, &calculation.CreatedAt); err != nil {
			return nil, fmt.Errorf("%v: %w", component, err)
		}
		calculations = append(calculations, &calculation)
	}

	if len(calculations) == 0 {
		return nil, domain.ErrNoCalculationsFound
	}

	return calculations, nil
}
