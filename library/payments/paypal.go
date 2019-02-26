package payments

import (
	"os"

	"github.com/logpacker/PayPal-Go-SDK"
)

type PayPal struct {
	*Detail
}

func (m *PayPal) Direct() error {
	if c, err := paypalsdk.NewClient(m.ClientID, m.SecretID, paypalsdk.APIBaseSandBox); err != nil {
		return err
	} else {
		c.SetLog(os.Stdout)
		if _, err = c.GetAccessToken(); err != nil {
			return err
		} else {
			payment := paypalsdk.Payment{
				Intent: m.PaymentIntent,
				Payer: &paypalsdk.Payer{
					PaymentMethod: m.PaymentMethod,
					FundingInstruments: []paypalsdk.FundingInstrument{{
						CreditCard: &paypalsdk.CreditCard{
							Number:      m.PaymentMethodCardNumber,
							Type:        m.PaymentMethodCardType,
							ExpireMonth: m.PaymentMethodCardExpireMonth,
							ExpireYear:  m.PaymentMethodCardExpireYear,
							CVV2:        m.PaymentMethodCardCVV2,
							FirstName:   m.PayerFirstName,
							LastName:    m.PayerLastName,
						},
					}},
				},
				Transactions: []paypalsdk.Transaction{{
					Amount: &paypalsdk.Amount{
						Currency: m.Currency,
						Total:    m.Total,
					},
					Description: m.PaymentDescription,
				}},
				RedirectURLs: &paypalsdk.RedirectURLs{
					ReturnURL: m.ReturnURL,
					CancelURL: m.CancelURL,
				},
			}
			_, err = c.CreatePayment(payment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *PayPal) CreditCardPayment() error {
	m.PaymentMethod = "credit_card"
	if c, err := paypalsdk.NewClient(m.ClientID, m.SecretID, paypalsdk.APIBaseSandBox); err != nil {
		return err
	} else {
		c.SetLog(os.Stdout)
		if _, err = c.GetAccessToken(); err != nil {
			return err
		} else {
			payment := paypalsdk.Payment{
				Intent: m.PaymentIntent,
				Payer: &paypalsdk.Payer{
					PaymentMethod: m.PaymentMethod,
					FundingInstruments: []paypalsdk.FundingInstrument{{
						CreditCard: &paypalsdk.CreditCard{
							Number:      m.PaymentMethodCardNumber,
							Type:        m.PaymentMethodCardType,
							ExpireMonth: m.PaymentMethodCardExpireMonth,
							ExpireYear:  m.PaymentMethodCardExpireYear,
							CVV2:        m.PaymentMethodCardCVV2,
							FirstName:   m.PayerFirstName,
							LastName:    m.PayerLastName,
						},
					}},
				},
				Transactions: []paypalsdk.Transaction{{
					Amount: &paypalsdk.Amount{
						Currency: m.Currency,
						Total:    m.Total,
					},
					Description: m.PaymentDescription,
				}},
				RedirectURLs: &paypalsdk.RedirectURLs{
					ReturnURL: m.ReturnURL,
					CancelURL: m.CancelURL,
				},
			}
			_, err = c.CreatePayment(payment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
