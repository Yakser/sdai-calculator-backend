package domain

import "time"

type Calculation struct {
	ID                         int64     `json:"id" db:"id"`
	UserID                     int64     `json:"user_id" db:"user_id"`
	PainfulJoints              int64     `json:"painful_joints" db:"painful_joints"`
	SwollenJoints              int64     `json:"swollen_joints" db:"swollen_joints"`
	PhysicalActivityAssessment int64     `json:"physician_activity_assessment" db:"physician_activity_assessment"`
	PatientActivityAssessment  int64     `json:"patient_activity_assessment" db:"patient_activity_assessment"`
	Crp                        float64   `json:"crp" db:"crp"`
	SdaiIndex                  string    `json:"sdai_index" db:"sdai_index"`
	CreatedAt                  time.Time `json:"created_at" db:"created_at"`
}
