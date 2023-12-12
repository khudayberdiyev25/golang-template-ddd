package app

import "github.com/google/uuid"

type NewCard struct {
	CustomerID uuid.UUID
	CardNumber string
	Expiry     string
	CardType   string
}
