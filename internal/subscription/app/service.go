package app

import (
	"errors"
	"github.com/google/uuid"
	customerdomain "gitlab.iman.uz/imandev/bnpl_payment/internal/subscription/domain/customer"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/subscription/domain/customer_card"
	err_pkg "gitlab.iman.uz/imandev/common_package/pkg/errors"
	"golang.org/x/net/context"
)

type SubscriptionService struct {
	cardRepo     customer_card.Repository
	cardProvider customer_card.Provider
	customerRepo customerdomain.Repository
	factory      customer_card.Factory
}

func NewSubscriptionService(
	cardRepo customer_card.Repository,
	cardProvider customer_card.Provider,
	customerRepo customerdomain.Repository,
	factory customer_card.Factory,
) *SubscriptionService {
	return &SubscriptionService{
		cardRepo:     cardRepo,
		cardProvider: cardProvider,
		customerRepo: customerRepo,
		factory:      factory,
	}
}

func (s *SubscriptionService) AddUnverifiedCard(ctx context.Context, newCard *NewCard) (*customer_card.Card, error) {
	card := s.factory.NewUnverifiedCard(
		newCard.CustomerID,
		newCard.CardNumber,
		newCard.Expiry,
		newCard.CardType,
	)

	c, err := s.getCardByCustomerAndNumber(ctx, newCard.CustomerID, newCard.CardNumber)
	if err != nil && !errors.Is(err, customer_card.ErrCardNotFound) {
		return nil, err
	}

	if c != nil && c.IsActive() {
		return nil, errors.New("card exists")
	}

	cardToken, err := s.cardProvider.CreateCard(ctx, card.CardNumber(), card.Expiry())
	if err != nil {
		return nil, err_pkg.NewSlugError(err.Error(), "cardProvider-unable-to-create-card")
	}
	card.SetToken(cardToken)

	if c != nil {
		card.SetGuid(c.Guid())

		if err = s.cardRepo.Update(ctx, card); err != nil {
			return nil, err_pkg.NewSlugError(err.Error(), "unable-to-update-card")
		}
	} else {
		if err = s.cardRepo.Create(ctx, card); err != nil {
			return nil, err_pkg.NewSlugError(err.Error(), "unable-to-create-card")
		}
	}

	return card, nil
}

func (s *SubscriptionService) VerifyCard(ctx context.Context, otpToken, otpCode string) error {
	cardGuid, err := uuid.Parse(otpToken)
	if err != nil {
		return err
	}

	card, err := s.cardRepo.FindOneByGuid(ctx, cardGuid)
	if err != nil {
		return err_pkg.NewSlugError(err.Error(), "unable-to-find-card-by-guid")
	}

	cardModel, err := s.cardProvider.VerifyCard(ctx, otpCode, card.Token())
	if err != nil {
		return err
	}

	customer, err := s.customerRepo.FindOneByGuid(ctx, card.CustomerID())
	if err != nil {
		return err_pkg.NewSlugError(err.Error(), "unable-to-find-customer-by-guid")
	}

	if err := card.VerifyCardOwnership(customer, cardModel); err != nil {
		card.MarkStatusActive()

		if err := s.cardRepo.Update(ctx, card); err != nil {
			return err_pkg.NewSlugError(err.Error(), "unable-to-update-card")
		}
	}

	return nil
}

func (s *SubscriptionService) getCardByCustomerAndNumber(ctx context.Context, customerID uuid.UUID, cardNumber string) (*customer_card.Card, error) {
	m := customer_card.ReadModel{}
	m.SetCustomerID(customerID)
	m.SetCardNumber(cardNumber)

	return s.cardRepo.FindOne(ctx, &m)
}
