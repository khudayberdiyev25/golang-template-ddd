package customer_card

import (
	"github.com/google/uuid"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/subscription/domain/customer"
	"time"
)

const (
	CardStatusActive      = "active"
	CardStatusDeactivated = "deactivated"
	CardStatusNotVerified = "not_verified"
	CardProviderPayme     = "payme_subscribe"
	CardProviderAlifPay   = "alifpay"
)

type Card struct {
	guid       uuid.UUID
	customerID uuid.UUID
	cardNumber string
	expiry     string
	cardType   string
	token      string
	status     string
	maskedPan  string
	owner      string
	own        bool
	provider   string
	createdAt  time.Time
	updatedAt  time.Time
}

func (c *Card) Guid() uuid.UUID {
	return c.guid
}

func (c *Card) SetGuid(guid uuid.UUID) {
	c.guid = guid
}

func (c *Card) CustomerID() uuid.UUID {
	return c.customerID
}

func (c *Card) CardNumber() string {
	return c.cardNumber
}

func (c *Card) Expiry() string {
	return c.expiry
}

func (c *Card) CardType() string {
	return c.cardType
}

func (c *Card) Token() string {
	return c.token
}

func (c *Card) SetToken(token string) {
	c.token = token
}

func (c *Card) Status() string {
	return c.status
}

func (c *Card) MaskedPan() string {
	return c.maskedPan
}

func (c *Card) Owner() string {
	return c.owner
}

func (c *Card) Own() bool {
	return c.own
}

func (c *Card) Provider() string {
	return c.provider
}

func (c *Card) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Card) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Card) IsActive() bool {
	return c.status == CardStatusActive && c.token != ""
}

func (c *Card) MarkStatusActive() {
	c.status = CardStatusActive
}

func (c *Card) VerifyCardOwnership(customer *customer.Customer, model *CardModel) error {
	c.maskedPan = model.MaskedPan
	c.owner = model.HolderName

	ok, _ := customer.MatchCardHolderName(c.Owner())
	if ok {
		c.own = true
	}

	return nil
}
