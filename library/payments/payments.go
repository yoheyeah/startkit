package payments

type (
	PaymentType string
)

const (
	PayPalPayment    PaymentType = "PayPal"
	AliPayPayment    PaymentType = "AliPay"
	WeChatPayPayment PaymentType = "WeChatPay"
	UnionPayPayment  PaymentType = "UnionPay"
)

type Detail struct {
	GatewayType      string
	IsProduction     bool
	Paymentdetail    Paymentdetail
	PayPalImplDetail PayPalImplDetail
	AliPayImplDetail AliPayImplDetail
}

type Paymentdetail struct {
	Subject           string
	ClientName        string
	ClientEmail       string
	ProductName       string
	Note              string
	Total             string
	DepositPercentage string
	Price             string
	Currency          string
}

type Payment interface {
	ReceivePayment() ([]byte, error)
	Refund() ([]byte, error)
	CreateOrder() (resp interface{}, err error)
	GetOrder(orderID string) (resp interface{}, err error)
	AuthorizeOrder(orderID string) (resp interface{}, err error)
	CaptureOrder(orderID string, isFinish bool) (resp interface{}, err error)
}

func GatewayType(p string, detail *Detail) Payment {
	switch PaymentType(p) {
	case PayPalPayment:
		return &PayPal{Detail: detail}
	case AliPayPayment:
		return &AliPay{Detail: detail}
	case WeChatPayPayment:
		return &WeChatPay{Detail: detail}
	case UnionPayPayment:
		return &UnionPay{Detail: detail}
	default:
		return &PayPal{Detail: detail}
	}
}

func ReceivePayment(detail Detail) ([]byte, error) {
	p := GatewayType(detail.GatewayType, &detail)
	return p.ReceivePayment()
}

func Refund(detail Detail) ([]byte, error) {
	p := GatewayType(detail.GatewayType, &detail)
	return p.Refund()
}

func CreateOrder(detail Detail) (resp interface{}, err error) {
	p := GatewayType(detail.GatewayType, &detail)
	return p.CreateOrder()
}

func GetOrder(orderID string, detail Detail) (resp interface{}, err error) {
	p := GatewayType(detail.GatewayType, &detail)
	return p.GetOrder(orderID)
}

func AuthorizeOrder(orderID string, detail Detail) (resp interface{}, err error) {
	p := GatewayType(detail.GatewayType, &detail)
	return p.AuthorizeOrder(orderID)
}

func CaptureOrder(orderID string, isFinish bool, detail Detail) (resp interface{}, err error) {
	p := GatewayType(detail.GatewayType, &detail)
	return p.CaptureOrder(orderID, isFinish)
}
