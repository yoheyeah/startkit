package email

import (
	"github.com/matcornic/hermes"
)

const (
	welcome        = "Welcome"
	information    = "Information"
	resetPwd       = "ResetPwd"
	activateAcount = "ActivateAcount"
	deleteAccount  = "DeleteAccount"
)

type Template interface {
	Email() hermes.Email
	Name(string) string
}

type Header struct {
	Type, Sender, Topic, Link, Logo string
	Receivers                       []string
}

func EmailInText() {

}

func EmailInHTML() {

}

func Type(t string) Template {
	switch t {
	case welcome:
		return Welcome{}
	case information:
		return nil
	case resetPwd:
		return nil
	case activateAcount:
		return nil
	case deleteAccount:
		return nil
	}
}

func EmailTemplate(header *Header) {
	m := Type(header.Type)
	return m.Email()
}
