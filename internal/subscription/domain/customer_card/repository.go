package customer_card

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type Repository interface {
	FindOneByGuid(ctx context.Context, guid uuid.UUID) (*Card, error)
	FindOne(ctx context.Context, model *ReadModel) (*Card, error)
	FindAll(ctx context.Context, model *ReadModel) ([]*Card, error)
	Create(ctx context.Context, card *Card) error
	Update(ctx context.Context, card *Card) error
}

const (
	ErrCardNotFound = Err("card not found")
)

type Err string

func (e Err) Error() string {
	return string(e)
}
