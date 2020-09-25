package payments

import (
	"encoding/json"

	"github.com/fatih/structs"
	paypalsdk "github.com/logpacker/PayPal-Go-SDK"
	paypalsdkv1 "github.com/logpacker/PayPal-Go-SDK-1.1.4"
)

type PayPalImplDetail struct {
	ClientID                     string
	SecretID                     string
	PaymentIntent                string
	PaymentMethod                string
	PaymentMethodCardType        string
	PaymentMethodCardNumber      string
	PaymentMethodCardExpireMonth string
	PaymentMethodCardExpireYear  string
	PaymentMethodCardCVV2        string
	PayerFirstName               string
	PayerLastName                string
	Total                        string
	Currency                     string
	ReturnURL                    string
	RedirectURL                  string
	CancelURL                    string
	PaymentDescription           string
}

type PayPal struct {
	*Detail
}

func (m *PayPal) ReceivePayment() (data []byte, err error) {
	client, err := m.NewClientv1()
	if err != nil {
		return
	}
	_, err = client.GetAccessToken()
	if err != nil {
		return
	}
	paymentResult, err := client.CreateDirectPaypalPayment(
		paypalsdkv1.Amount{
			Total:    m.Paymentdetail.Price,
			Currency: m.Paymentdetail.Currency,
		},
		m.PayPalImplDetail.RedirectURL,
		m.PayPalImplDetail.CancelURL,
		m.PayPalImplDetail.PaymentDescription,
	)
	if err != nil {
		return
	}
	payment, err := client.GetPayment(paymentResult.ID)
	if err != nil {
		return
	}
	executeResult, err := client.ExecuteApprovedPayment(paymentResult.ID, payment.Payer.PayerInfo.PayerID)
	if err != nil {
		return
	}
	data, err = json.Marshal(structs.Map(executeResult))
	return
}

func (m *PayPal) Refund() (data []byte, err error) {
	c, err := m.NewClientv1()
	if err != nil {
		return nil, err
	}
	if _, err = c.GetAccessToken(); err != nil {
		return nil, err
	}
	c.CreateSinglePayout(paypalsdkv1.Payout{
		SenderBatchHeader: &paypalsdkv1.SenderBatchHeader{
			EmailSubject: m.Paymentdetail.Subject,
		},
		Items: []paypalsdkv1.PayoutItem{
			paypalsdkv1.PayoutItem{
				RecipientType: "EMAIL",
				Receiver:      m.Paymentdetail.ClientEmail,
				Amount: &paypalsdkv1.AmountPayout{
					Value:    m.Paymentdetail.Price,
					Currency: m.Paymentdetail.Currency,
				},
				Note: m.Paymentdetail.Note,
			},
		},
	})
	return
}

func (m *PayPal) GetOrder(orderID string) (resp interface{}, err error) {
	client, err := m.NewClientv2()
	if err != nil {
		return
	}
	_, err = client.GetAccessToken()
	if err != nil {
		return
	}
	order, err := client.GetOrder(orderID)
	if err != nil {
		return
	}
	return order, err
}

func (m *PayPal) CreateOrder() (order interface{}, err error) {
	client, err := m.NewClientv2()
	if err != nil {
		return
	}
	_, err = client.GetAccessToken()
	if err != nil {
		return
	}
	createCheckout, err := client.CreateOrder(m.PayPalImplDetail.PaymentIntent, &paypalsdk.PurchaseUnitAmount{
		Currency: m.Paymentdetail.Currency,
		Value:    m.Paymentdetail.Price,
	})
	if err != nil {
		return "", err
	}
	return createCheckout, nil
}

func (m *PayPal) AuthorizeOrder(orderID string) (resp interface{}, err error) {
	client, err := m.NewClientv2()
	if err != nil {
		return
	}
	_, err = client.GetAccessToken()
	if err != nil {
		return
	}
	auth, err := client.AuthorizeOrder(orderID, m.PayPalImplDetail.PaymentIntent, &paypalsdk.PurchaseUnitAmount{
		Currency: m.Paymentdetail.Currency,
		Value:    m.Paymentdetail.Price,
	})
	if err != nil {
		return "", err
	}
	return auth, nil
}

func (m *PayPal) CaptureOrder(orderID string, isFinish bool) (resp interface{}, err error) {
	client, err := m.NewClientv2()
	if err != nil {
		return
	}
	_, err = client.GetAccessToken()
	if err != nil {
		return
	}
	auth, err := client.CaptureOrder(orderID, &paypalsdk.PurchaseUnitAmount{
		Currency: m.Paymentdetail.Currency,
		Value:    m.Paymentdetail.Price,
	},
		isFinish,
		&paypalsdk.Currency{
			Currency: m.Paymentdetail.Currency,
			Value:    m.Paymentdetail.Price,
		})
	if err != nil {
		return "", err
	}
	return auth, nil
}

func (m *PayPal) Direct() error {
	return nil
}

func (m *PayPal) CreditCardPayment() error {
	m.PayPalImplDetail.PaymentMethod = "credit_card"
	if c, err := m.NewClientv1(); err != nil {
		return err
	} else {
		if _, err = c.GetAccessToken(); err != nil {
			return err
		}
	}
	return nil
}

func (m *PayPal) NewClientv1() (client *paypalsdkv1.Client, err error) {
	if m.IsProduction {
		return paypalsdkv1.NewClient(m.PayPalImplDetail.ClientID, m.PayPalImplDetail.SecretID, paypalsdkv1.APIBaseLive)
	}
	return paypalsdkv1.NewClient(m.PayPalImplDetail.ClientID, m.PayPalImplDetail.SecretID, paypalsdkv1.APIBaseSandBox)
}

func (m *PayPal) Template(c *paypalsdkv1.Client) (err error) {

	return
}

func (m *PayPal) AccessTokenv1(c *paypalsdkv1.Client) (token *paypalsdkv1.TokenResponse, err error) {
	token, err = c.GetAccessToken()
	return
}

func (m *PayPal) GetAuthorizationByIDv1(c *paypalsdkv1.Client, authID string) (auth *paypalsdkv1.Authorization, err error) {
	auth, err = c.GetAuthorization(authID)
	return
}

func (m *PayPal) CaptureAuthorizationv1(c *paypalsdkv1.Client, authID string) (err error) {
	// capture, err := c.CaptureAuthorization(authID, &paypalsdkv1.Amount{Total: "7.00", Currency: "USD"}, true)
	return
}

func (m *PayPal) NewClientv2() (client *paypalsdk.Client, err error) {
	if m.IsProduction {
		return paypalsdk.NewClient(m.PayPalImplDetail.ClientID, m.PayPalImplDetail.SecretID, paypalsdk.APIBaseLive)
	}
	return paypalsdk.NewClient(m.PayPalImplDetail.ClientID, m.PayPalImplDetail.SecretID, paypalsdk.APIBaseSandBox)
}

func (m *PayPal) AccessTokenv2(c *paypalsdk.Client) (token *paypalsdk.TokenResponse, err error) {
	token, err = c.GetAccessToken()
	return
}

func (m *PayPal) GetAuthorizationByIDv2(c *paypalsdk.Client, authID string) (auth *paypalsdk.Authorization, err error) {
	auth, err = c.GetAuthorization(authID)
	return
}

func (m *PayPal) CaptureAuthorizationv2(c *paypalsdk.Client, authID string) (err error) {
	// capture, err := c.CaptureAuthorization(authID, &paypalsdkv1.Amount{Total: "7.00", Currency: "USD"}, true)
	return
}
