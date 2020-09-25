package payments

type UnionPay struct {
	*Detail
}

func (m *UnionPay) ReceivePayment() (data []byte, err error) {
	return
}

func (m *UnionPay) Refund() (data []byte, err error) {
	return
}

func (m *UnionPay) CreateOrder() (resp interface{}, err error) {
	return
}

func (m *UnionPay) GetOrder(orderID string) (resp interface{}, err error) {
	return
}

func (m *UnionPay) AuthorizeOrder(orderID string) (resp interface{}, err error) {
	return
}

func (m *UnionPay) CaptureOrder(orderID string, isFinish bool) (resp interface{}, err error) {
	return
}
