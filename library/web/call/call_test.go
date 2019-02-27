package call

import (
	"testing"
)

func TestRequest_BuildURL(t *testing.T) {
	type fields struct {
		IsHTTPS     bool
		Method      string
		UserInfo    string
		Host        string
		Port        int
		Params      []string
		Querys      map[string]string
		Fragments   []string
		RequestBody interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantUrl string
	}{
		// TODO: Add test cases.
		{
			name: "functionality 1",
			fields: fields{
				IsHTTPS:     false,
				Method:      "GET",
				UserInfo:    "",
				Host:        "testfunctionality.com",
				Port:        0,
				Params:      []string{"param1", "param2"},
				Querys:      map[string]string{"query1": "value1", "query2": "value2"},
				Fragments:   []string{},
				RequestBody: nil,
			},
			wantUrl: "",
		},
		{
			name: "functionality 2",
			fields: fields{
				IsHTTPS:     true,
				Method:      "GET",
				UserInfo:    "testuser1",
				Host:        "testfunctionality.com",
				Port:        1000,
				Params:      []string{"param1", "param2"},
				Querys:      map[string]string{"query1": "value1", "query2": "value2"},
				Fragments:   []string{"fragment1", "fragment2"},
				RequestBody: nil,
			},
			wantUrl: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				IsHTTPS:     tt.fields.IsHTTPS,
				Method:      tt.fields.Method,
				UserInfo:    tt.fields.UserInfo,
				Host:        tt.fields.Host,
				Port:        tt.fields.Port,
				Params:      tt.fields.Params,
				Querys:      tt.fields.Querys,
				Fragments:   tt.fields.Fragments,
				RequestBody: tt.fields.RequestBody,
			}
			if gotUrl := r.BuildURL(); gotUrl != tt.wantUrl {
				t.Errorf("Request.BuildURL() = %v, want %v", gotUrl, tt.wantUrl)
			}
		})
	}
}

func TestHTTPRequest(t *testing.T) {
	type args struct {
		req *Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test functionality",
			args: args{
				req: &Request{
					IsHTTPS:     false,
					Method:      "GET",
					UserInfo:    "",
					Host:        "www.66ip.cn/mo.php",
					Port:        0,
					Params:      []string{},
					Querys:      map[string]string{"tqsl": "100"},
					Fragments:   []string{},
					RequestBody: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HTTPRequest(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("HTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
