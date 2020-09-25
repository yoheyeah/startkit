package email

import (
	"testing"
)

func TestFor(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
		{
			"a",
			"a,b,c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := For(); got != tt.want {
				t.Errorf("For() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetter_SendMail(t *testing.T) {
	type fields struct {
		Provider  string
		Host      string
		Port      string
		User      string
		Password  string
		Type      string
		Sender    string
		Subject   string
		Topic     string
		Link      string
		Logo      string
		Content   string
		Receivers []string
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "aaaaaa",
			fields: fields{
				Host:      "smtp.office365.com",
				Port:      "587",
				User:      "starhero@live.hk",
				Password:  "SH138161888",
				Sender:    "starhero@live.hk",
				Subject:   "aaaaaa",
				Topic:     "aaaaaa",
				Receivers: []string{"harrisin2037@gmail.com"},
			},
			args: args{
				email: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Setter{
				Provider:  tt.fields.Provider,
				Host:      tt.fields.Host,
				Port:      tt.fields.Port,
				User:      tt.fields.User,
				Password:  tt.fields.Password,
				Type:      tt.fields.Type,
				Sender:    tt.fields.Sender,
				Subject:   tt.fields.Subject,
				Topic:     tt.fields.Topic,
				Link:      tt.fields.Link,
				Logo:      tt.fields.Logo,
				Content:   tt.fields.Content,
				Receivers: tt.fields.Receivers,
			}
			if err := m.SendMail(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("Setter.SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
