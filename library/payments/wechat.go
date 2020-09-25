package payments

type WeChatPay struct {
	*Detail
}

func (m *WeChatPay) ReceivePayment() (data []byte, err error) {
	return
}

func (m *WeChatPay) Refund() (data []byte, err error) {
	return
}

func (m *WeChatPay) CreateOrder() (resp interface{}, err error) {
	return
}

func (m *WeChatPay) GetOrder(orderID string) (resp interface{}, err error) {
	return
}

func (m *WeChatPay) AuthorizeOrder(orderID string) (resp interface{}, err error) {
	return
}

func (m *WeChatPay) CaptureOrder(orderID string, isFinish bool) (resp interface{}, err error) {
	return
}
