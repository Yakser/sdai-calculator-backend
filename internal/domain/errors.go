package domain

import "errors"

var (
	ErrCalculationAlreadyExists = errors.New("calculation already exists")
	ErrNoCalculationsFound      = errors.New("no calculations found")
)
