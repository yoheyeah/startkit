package email

import (
	"github.com/matcornic/hermes"
)

type Welcome struct {
	ReceiverName  string
	Introductions []string
	Lists         map[string]string
	Actions       []struct {
		Instructions string
		ButtonText   string
		ButtonLink   string
	}
	Ending []string
}

func (m *Welcome) Default() Welcome {
	m = &Welcome{
		
	}
	return m
}

func (m *Welcome) Name(name string) string {
	return name
}

func (m *Welcome) Email() hermes.Email {
	var (
		dict    = []hermes.Entry{}
		actions = []hermes.Action{}
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
	return hermes.Email{
		Body: hermes.Body{
			Name:       m.ReceiverName,
			Intros:     Introductions,
			Dictionary: dict,
			Actions:    actions,
			Outros:     m.Ending,
		},
	}
}
