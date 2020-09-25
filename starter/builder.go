package starter

import (
	"errors"
	"log"
	"startkit/library/files"
	"startkit/library/systems"
)

type Builder interface {
	Builder(*Content) error
}

var (
	defaultConfigFiles = []string{
		"default_setting.ini",
		"default_setting.json",
	}
)

var (
	_ Builder = &Config{}
	_ Builder = &Logger{}
	_ Builder = &App{}
	_ Builder = &Server{}
	_ Builder = &Mysql{}
	_ Builder = &Mongo{}
	_ Builder = &Influx{}
	_ Builder = &Redis{}
)

var (
	NotExistConfigFileError = errors.New("Config file provided not exist")
)

func DefaultBuilder() (content *Content) {
	for i := 0; i < len(defaultConfigFiles); i++ {
		if !systems.IsNotExist(defaultConfigFiles[i]) {
			content = &Content{ConfigFile: defaultConfigFiles[i]}
		}
	}
	files.BindFileToObj(content.ConfigFile, content)
	for i := 0; i < len(content.App.InUseService); i++ {
		err := content.Builder(content.GetFieldOfStructPointer(content.App.InUseService[i]))
		if err != nil {
			err = errors.New(content.App.InUseService[i] + ":" + err.Error())
			content.Errors = append(content.Errors, err)
		} else {
			log.Println(content.App.InUseService[i] + " started")
		}
	}
	return
}

func CustomBuilder(file string) (content *Content) {
	if !systems.IsNotExist(file) {
		content = &Content{ConfigFile: file}
	}
	if err := files.BindFileToObj(content.ConfigFile, content); err != nil {
		panic(err)
	}
	for i := 0; i < len(content.App.InUseService); i++ {
		err := content.Builder(content.GetFieldOfStructPointer(content.App.InUseService[i]))
		if err != nil {
			err = errors.New(content.App.InUseService[i] + ":" + err.Error())
			content.Errors = append(content.Errors, err)
		} else {
			log.Println(content.App.InUseService[i] + " started")
		}
	}
	return
}

func (m *Content) Builder(b Builder) error {
	return b.Builder(m)
}
