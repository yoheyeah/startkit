package email

import (
	"github.com/matcornic/hermes"
)

type Welcome struct {
	CompanyName     string
	CompanyLink     string
	CompanyLogoLink string
	ReceiverName    string
	Introductions   []string
	Lists           map[string]string
	Actions         []struct {
		Instructions string
		ButtonText   string
		ButtonLink   string
	}
	Ending []string
}

func (m *Welcome) Email() (emailText string) {
	var (
		err     error
		dict    = []hermes.Entry{}
		actions = []hermes.Action{}
		header  = hermes.Hermes{
			Theme: new(hermes.Default),
			Product: hermes.Product{
				Name: m.CompanyName,
				Link: m.CompanyLink,
				Logo: m.CompanyLogoLink,
			},
		}
	)
	for k, v := range m.Lists {
		dict = append(dict, hermes.Entry{
			Key:   k,
			Value: v,
		})
	}
	if count := len(m.Actions); count > 0 {
		for i := 0; i < count; i++ {
			actions = append(actions, hermes.Action{
				Instructions: m.Actions[i].Instructions,
				Button: hermes.Button{
					Text: m.Actions[i].ButtonText,
					Link: m.Actions[i].ButtonLink,
				},
			})
		}
	}
	email := hermes.Email{
		Body: hermes.Body{
			Name:       m.ReceiverName,
			Intros:     m.Introductions,
			Dictionary: dict,
			Actions:    actions,
			Outros:     m.Ending,
		},
	}
	if emailText, err = header.GeneratePlainText(email); err != nil {
		return ""
	}
	return emailText
}
