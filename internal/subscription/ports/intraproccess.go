package ports

import (
	"github.com/google/uuid"
	gp "gitlab.iman.uz/imandev/bnpl_contracts/genproto/bnpl_payment_service"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/subscription/app"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/subscription/domain/customer_card"
	"golang.org/x/net/context"
)

type CardInterface struct {
	service app.SubscriptionService
}

func NewCardInterface(service app.SubscriptionService) *CardInterface {
	return &CardInterface{service: service}
}

func (s *CardInterface) AddCard(ctx context.Context, in *gp.AddCardRequest) (*gp.AddCardResponse, error) {
	customerId, err := uuid.Parse(in.CustomerId)
	if err != nil {
		return nil, err
	}

	nC := app.NewCard{
		CustomerID: customerId,
		CardNumber: in.CardNumber,
		Expiry:     in.Expiry,
		CardType:   in.CardType,
	}

	c, err := s.service.AddUnverifiedCard(ctx, &nC)
	if err != nil {
		return nil, err
	}

	resp := gp.AddCardResponse{
		OtpToken: c.Guid().String(),
		Card:     s.cardFromDomain(c),
	}

	return &resp, nil
}

func (s *CardInterface) cardFromDomain(domainCard *customer_card.Card) *gp.Card {
	c := &gp.Card{
		Guid:       domainCard.Guid().String(),
		CustomerId: domainCard.CustomerID().String(),
		CardNumber: domainCard.CardNumber(),
		Expiry:     domainCard.Expiry(),
		CardType:   domainCard.CardType(),
		Status:     domainCard.Status(),
		Owner:      domainCard.Owner(),
		Own:        domainCard.Own(),
		MaskedPan:  domainCard.MaskedPan(),
		CreatedAt:  domainCard.CreatedAt().String(),
		UpdatedAt:  domainCard.UpdatedAt().String(),
	}

	if domainCard.IsActive() {
		c.IsSubscribed = true
		c.IsConfirm = false
	} else {
		c.IsConfirm = true
	}

	return c
}
