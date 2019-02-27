package starter

import (
	"log"
	"startkit/library/debugs"
	"startkit/library/files"
	"startkit/library/systems"
	"strconv"
)

type Config struct {
	File   string
	App    App
	Server Server
	Mysql  Mysql
	Mongo  Mongo
}

func (m *Config) Builder(c *Content) error {
	files.BindFileToObj(m.File, m)
	version, err := systems.GetMinimumVersion(m.App.MinimumGoVersion)
	if err != nil {
		log.Fatalln(err)
		c.Errors = append(c.Errors, err)
	}
	if v, err := systems.GetMinimumVersion(""); v <= version {
		debugs.DebuggingPrint(`[WARNING] Now require Go version ` + strconv.Itoa(int(v)) + ` or later. `)
	} else if err != nil {
		log.Fatalln(err)
		c.Errors = append(c.Errors, err)
	}
	debugs.DebuggingPrint(`[WARNING] Building an Config instance. `)
	return nil
}

func (m *Config) Starter(c *Content) error {
	return nil
}

func (m *Config) Router(s *Server) {

}
