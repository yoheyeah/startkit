package email

import (
	"io"
	"net/smtp"

	"github.com/domodwyer/mailyak"
)

type loginAuth struct {
	username, password string
}

func genLoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}

func (m *Setter) Send(email string) (err error) {
	var (
		auth = smtp.PlainAuth("", m.User, m.Password, m.Host)
		mail = &mailyak.MailYak{}
	)
	if m.Provider == "hotmail" {
		return m.SendMail(email)
	} else {
		mail = mailyak.New(m.Host+":"+m.Port, auth)
	}
	mail.To(m.Receivers...)
	mail.From(m.User)
	mail.FromName(m.Sender)
	mail.Subject(m.Subject)
	if m.ReplyTo != "" {
		mail.ReplyTo(m.ReplyTo)
	}
	mail.AddHeader("X-TOTALLY-NOT-A-SCAM", "true")
	if _, err = io.WriteString(mail.HTML(), email); err != nil {
		return
	}
	if err = mail.Send(); err != nil {
		return
	}
	return
}

func (m *Setter) SendMail(email string) error {
	auth := genLoginAuth(m.User, m.Password)
	target := ""
	if count := len(m.Receivers); count > 0 {
		target = m.Receivers[0]
		if count > 1 {
			for i := 1; i < count; i++ {
				target = target + "," + m.Receivers[i]
			}
		}
	}
	contentType := "Content-Type: text/plain" + "; charset=UTF-8"
	msg := []byte("To: " + target +
		"\r\nFrom: " + m.User +
		"\r\nSubject: " + m.Subject +
		"\r\n" + contentType + "\r\n\r\n" +
		email)
	err := smtp.SendMail(m.Host+":"+m.Port, auth, m.User, m.Receivers, msg)
	if err != nil {
		return err
	}
	return nil
}

func For() string {
	a := []string{
		"a", "b", "c",
	}
	target := ""
	if count := len(a); count > 0 {
		target = a[0]
		if count > 1 {
			for i := 1; i < count; i++ {
				target = target + "," + a[i]
			}
		}
	}
	return target
}
