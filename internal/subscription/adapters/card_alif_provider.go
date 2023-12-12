package adapters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/common/config"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/subscription/domain/provider"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"log"
	"time"
)

type alifpayProvider struct {
	logger     *zap.Logger
	conf       *config.Config
	httpClient *resty.Client
}

func NewCardProviderALIF(conf *config.Config, logger *zap.Logger) (provider.Provider, error) {
	timeout, err := time.ParseDuration("30s")
	if err != nil {
		return nil, fmt.Errorf("error during parse duration for alif_pay timeout : %w", err)
	}

	httpClient := resty.New().
		SetBaseURL(conf.ALIF.BaseURL).
		SetTimeout(timeout).
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			body, err := json.Marshal(r.Body)
			if err != nil {
				logger.Error("alif_pay before request error", zap.Error(err))
			}
			log.Println("alif_pay request", string(body))
			fmt.Println()
			return nil
		}).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			log.Println("alif_pay response", resp.String())
			fmt.Println()
			return nil
		}).
		SetHeader("Token", conf.ALIF.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	return &alifpayProvider{
		logger:     logger,
		conf:       conf,
		httpClient: httpClient,
	}, nil
}

func (p alifpayProvider) CreateCard(ctx context.Context, number, expire string) (cardToken string, err error) {
	body := map[string]any{
		"pan": number,
		"exp": expire,
	}

	httpResp, err := p.httpClient.R().
		SetBody(body).
		Post("/getCardToken")

	if err != nil {
		return "", err
	}

	if httpResp.StatusCode() != 200 {
		return "", errors.New("error response: " + httpResp.String())
	}

	var resp createCardResponse
	err = json.Unmarshal(httpResp.Body(), &resp)
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", errors.New(resp.Error.Message)
	}

	return resp.Token, nil
}

func (p alifpayProvider) VerifyCard(ctx context.Context, code, cardToken string) (*provider.CardModel, error) {
	body := map[string]any{
		"token": cardToken,
		"otp":   code,
	}

	httpResp, err := p.httpClient.R().
		SetBody(body).
		Post("/confirmCardToken")

	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode() != 200 {
		return nil, errors.New("error response: " + httpResp.String())
	}

	var resp verifyCardResponse
	err = json.Unmarshal(httpResp.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	model := provider.CardModel{
		MaskedPan:   resp.MaskedPan,
		BankName:    resp.BankName,
		HolderName:  resp.HolderName,
		Token:       resp.Token,
		MaskedPhone: resp.MaskedPhone,
	}

	return &model, nil
}

type createCardResponse struct {
	Token       string `json:"token"`
	MaskedPhone string `json:"masked_phone"`
	Error       *err   `json:"error"`
}

type verifyCardResponse struct {
	card
	Error *err `json:"error"`
}

type card struct {
	MaskedPan   string `json:"masked_pan"`
	BankName    string `json:"bank_name"`
	HolderName  string `json:"holder_name"`
	Token       string `json:"token"`
	MaskedPhone string `json:"masked_phone"`
}

type err struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
