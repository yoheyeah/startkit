package payments

// import (
// 	"github.com/smartwalle/alipay"
// )

type AliPayImplDetail struct {
	AppID       string
	PublicKey   string
	PrivateKey  string
	NotifyURL   string
	ReturnURL   string
	Subject     string
	OutTradeNo  string // unique invoice no
	TotalAmount string
	ProductCode string
}

type AliPay struct {
	*Detail
	*AliPayImplDetail
}

func (m *AliPay) ReceivePayment() (data []byte, err error) {
	return
}

func (m *AliPay) Refund() (data []byte, err error) {
	return
}

func (m *AliPay) CreateOrder() (resp interface{}, err error) {
	return
}

func (m *AliPay) GetOrder(orderID string) (resp interface{}, err error) {
	return
}

func (m *AliPay) AuthorizeOrder(orderID string) (resp interface{}, err error) {
	return
}

func (m *AliPay) CaptureOrder(orderID string, isFinish bool) (resp interface{}, err error) {
	return
}

func (m *AliPay) NewClient() {
	// var client = alipay.New(m.AppID, m.PublicKey, m.PrivateKey, m.IsProduction)
	// var p = alipay.AliPayTradePagePay{}
	// p.NotifyURL = "http://220.112.233.229:3000/alipay"
	// p.ReturnURL = "http://220.112.233.229:3000"
	// p.Subject = "修正了中文的 Bug"
	// p.OutTradeNo = "trade_no_20170623011121"
	// p.TotalAmount = "10.00"
	// p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	// url, err := client.TradePagePay(p)
	// if err != nil {

	// }
}
