package domain

import "time"

type Calculation struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	SdaiIndex string    `json:"sdai_index"`
	CreatedAt time.Time `json:"created_at"`
}
