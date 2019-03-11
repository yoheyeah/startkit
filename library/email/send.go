package email

import (
	"io"
	"net/smtp"

	"github.com/domodwyer/mailyak"
)

func (m *Setter) Send(email string) (err error) {
	var (
		mail = mailyak.New(m.Host+":"+m.Port, smtp.PlainAuth("", m.User, m.Password, m.Host))
	)
	mail.To(m.Receivers...)
	mail.From(m.User)
	mail.FromName(m.Sender)
	mail.Subject(m.Subject)
	mail.AddHeader("X-TOTALLY-NOT-A-SCAM", "true")
	if _, err = io.WriteString(mail.HTML(), email); err != nil {
		return
	}
	if err = mail.Send(); err != nil {
		return
	}
	return
}
