package call

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type GET struct {
	Request  *Request
	URL      string
	Response *http.Response
	RespStr  string
}

func (g *GET) BindRequest() error {
	g.Build()
	resp, err := http.Get(g.URL)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return err
		}
		return errors.New("Return With Status Code: " + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	} else {
		g.RespStr = string(body)
	}
	return nil
}

func (g *GET) Build() {
	g.URL = g.Request.BuildURL()
}
