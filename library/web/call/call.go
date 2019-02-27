package call

import (
	"strconv"
)

type Request struct {
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

type BindingRequest interface {
	BindRequest() error
}

func GetRequestType(req *Request) BindingRequest {
	switch req.Method {
	case "GET":
		return &GET{Request: req}
	default:
		return &GET{Request: req}
	}
}

// Api: the file name and the pointer of the struct variable
func HTTPRequest(req *Request) error {
	return GetRequestType(req).BindRequest()
}

func (r *Request) BuildURL() (url string) {
	if r.IsHTTPS {
		url = "https://"
	} else {
		url = "http://"
	}
	if r.UserInfo != "" {
		url = url + r.UserInfo + "@"
	}
	if r.Port != -1 && r.Port != 0 {
		url = url + r.Host + ":" + strconv.Itoa(r.Port)
	} else {
		url = url + r.Host
	}
	if count := len(r.Params); count > 0 {
		for i := 0; i < count; i++ {
			url = url + "/" + r.Params[i]
		}
	}
	if count := len(r.Querys); count > 0 {
		url = url + "?"
		for query, value := range r.Querys {
			url = url + query + "=" + value
			if count > 1 {
				url = url + "&"
			}
			count--
		}
	}
	if count := len(r.Fragments); count > 0 {
		for i := 0; i < count; i++ {
			url = url + "#" + r.Fragments[i]
		}
	}
	return
}
