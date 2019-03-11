package email

type Plain struct {
	Content string
}

func (m *Plain) Email() (emailText string) {
	emailText = m.Content
	return
}
