package customer_card

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type Factory struct {
}

func (f Factory) NewUnverifiedCard(
	customerID uuid.UUID,
	cardNumber string,
	expiry string,
	cardType string,
) *Card {
	expiry = f.sanitizeCardExpiry(expiry)
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")

	return &Card{
		guid:       uuid.New(),
		customerID: customerID,
		cardNumber: cardNumber,
		expiry:     expiry,
		cardType:   cardType,
		status:     CardStatusNotVerified,
		provider:   CardProviderAlifPay,
		createdAt:  time.Now().UTC(),
		updatedAt:  time.Now().UTC(),
	}
}

func (f Factory) UnmarshallCardFromDatabase(
	guid uuid.UUID,
	customerID uuid.UUID,
	cardNumber string,
	expiry string,
	cardType string,
	token string,
	status string,
	maskedPan string,
	owner string,
	own bool,
	provider string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Card, error) {

	return &Card{
		guid:       guid,
		customerID: customerID,
		cardNumber: cardNumber,
		expiry:     expiry,
		cardType:   cardType,
		token:      token,
		status:     status,
		maskedPan:  maskedPan,
		owner:      owner,
		own:        own,
		provider:   provider,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}, nil
}

func (f Factory) sanitizeCardExpiry(expiry string) string {
	switch expiry[0] {
	case 48, 49:
	default:
		expiry = expiry[2:] + expiry[:2]
	}

	return expiry
}
