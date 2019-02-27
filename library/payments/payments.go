package payments

type Detail struct {
	GatewayType                  string
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

type Payment interface {
	Direct() error
	CreditCardPayment() error
}

func GatewayType(p string, detail *Detail) Payment {
	switch p {
	case "PayPal":
		return &PayPal{Detail: detail}
	default:
		return &PayPal{Detail: detail}
	}
}

func Direct(detail Detail) error {
	p := GatewayType(detail.GatewayType, &detail)
	return p.Direct()
}

func CreditCardPayment(detail Detail) error {
	p := GatewayType(detail.GatewayType, &detail)
	return p.CreditCardPayment()
}
