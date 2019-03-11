package email

type Plain struct {
	*Setter
}

func (m *Plain) Email() (emailText string) {
	emailText = m.Content
	return
}
