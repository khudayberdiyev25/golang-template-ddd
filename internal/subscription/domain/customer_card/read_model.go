package customer_card

import (
	"github.com/google/uuid"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/common/pagination"
)

type ReadModel struct {
	customerID *uuid.UUID `json:"customer_id,omitempty"`
	cardNumber *string    `json:"card_number,omitempty"`
	pagination.Pagination
}

func NewReadModel() *ReadModel {
	return &ReadModel{
		customerID: nil,
		cardNumber: nil,
	}
}

func (r *ReadModel) CustomerID() (val uuid.UUID, ok bool) {
	if r.customerID == nil {
		return uuid.UUID{}, false
	}

	return *r.customerID, true
}

func (r *ReadModel) SetCustomerID(customerID uuid.UUID) {
	r.customerID = &customerID
}

func (r *ReadModel) CardNumber() (val string, ok bool) {
	if r.cardNumber == nil {
		return "", false
	}

	return *r.cardNumber, true
}

func (r *ReadModel) SetCardNumber(cardNumber string) {
	r.cardNumber = &cardNumber
}
