package ports

import (
	"context"
	"github.com/google/uuid"
	gp "gitlab.iman.uz/imandev/bnpl_contracts/genproto/bnpl_payment_service"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/subscription/app"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/subscription/domain/customer_card"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	service app.SubscriptionService
}

func NewGrpcServer(service app.SubscriptionService) gp.BNPLPaymentServer {
	return &GrpcServer{service: service}
}

func (h *GrpcServer) AddCard(ctx context.Context, in *gp.AddCardRequest) (*gp.AddCardResponse, error) {
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

	c, err := h.service.AddUnverifiedCard(ctx, &nC)
	if err != nil {
		return nil, err
	}

	resp := gp.AddCardResponse{
		OtpToken: c.Guid().String(),
		Card:     h.cardFromDomain(c),
	}

	return &resp, nil
}

func (h *GrpcServer) VerifyCard(ctx context.Context, in *gp.VerifyCardRequest) (*emptypb.Empty, error) {
	err := h.service.VerifyCard(ctx, in.GetOtpToken(), in.GetCode())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *GrpcServer) SubscribeCard(ctx context.Context, in *gp.SubscribeCardRequest) (*gp.SubscribeCardResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *GrpcServer) GetCard(ctx context.Context, in *gp.GetCardRequest) (*gp.GetCardResponse, error) {

	//TODO implement me
	panic("implement me")
}

func (h *GrpcServer) GetCardsList(ctx context.Context, in *gp.GetCardsListRequest) (*gp.GetCardsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *GrpcServer) RegisterTransaction(ctx context.Context, in *gp.RegisterTransactionRequest) (*gp.RegisterTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *GrpcServer) WithdrawAutoPayment(ctx context.Context, in *gp.WithdrawAutoPaymentRequest) (*gp.WithdrawAutoPaymentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *GrpcServer) WithdrawPayment(ctx context.Context, in *gp.WithdrawPaymentRequest) (*gp.WithdrawPaymentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *GrpcServer) PublishTransaction(ctx context.Context, in *gp.PublishTransactionRequest) (*gp.PublishTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *GrpcServer) SendFiscalCheck(ctx context.Context, in *gp.SendFiscalCheckRequest) (*gp.SendFiscalCheckResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *GrpcServer) Transactions(ctx context.Context, in *gp.TransactionsRequest) (*gp.TransactionsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *GrpcServer) cardFromDomain(domainCard *customer_card.Card) *gp.Card {
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
