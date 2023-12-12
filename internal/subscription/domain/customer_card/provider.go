package customer_card

import "golang.org/x/net/context"

type Provider interface {
	CreateCard(ctx context.Context, number, expire string) (cardToken string, err error)
	VerifyCard(ctx context.Context, code, cardToken string) (*CardModel, error)
}

type CardModel struct {
	MaskedPan   string
	BankName    string
	HolderName  string
	Token       string
	MaskedPhone string
}
