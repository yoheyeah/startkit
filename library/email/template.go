package email

const (
	welcome        = "Welcome"
	information    = "Information"
	resetPwd       = "ResetPwd"
	activateAcount = "ActivateAcount"
	deleteAccount  = "DeleteAccount"
)

type Template interface {
	Email() string
}

type Setter struct {
	Host, Port                               string
	User, Password                           string
	Type, Sender, Subject, Topic, Link, Logo string
	Receivers                                []string
}

func Type(t string) Template {
	switch t {
	case welcome:
		return &Welcome{}
	case information:
		return nil
	case resetPwd:
		return nil
	case activateAcount:
		return nil
	case deleteAccount:
		return nil
	default:
		return &Plain{}
	}
}

func EmailTemplate(setter *Setter) (err error) {
	return setter.Send(Type(setter.Type).Email())
}
