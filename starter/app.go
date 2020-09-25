package starter

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"runtime"
	"startkit/library/files"
	"startkit/library/systems"
	"startkit/library/times"
	"time"
)

type App struct {
	ModuleName   string
	ModuleID     int
	InUseService []string
	DebugMode    bool

	// basic setting
	ExecPath         string
	RootPath         string
	MinimumGoVersion string

	FileEncryptKey string

	// upload setting
	UploadPath      string
	UploadSizeLimit int
	UploadFileTypes []string

	// download setting
	DownloadPath      string
	DownloadSizeLimit int
	DownloadFileTypes []string

	ThumbNailPath string
	ThumbNailSize int

	// file location setting
	FileLocationShiftInterval int // Hourly

	// System variable
	MaxCPUThread int
	OSType       string
	IsWindows    bool
	IsLinux      bool
}

func (m *App) Builder(c *Content) error {
	var (
		start = []int{0, 0, 0, 0}
		err   error
	)
	m.IsLinux = false
	m.IsWindows = false
	m.ExecPath, err = files.ExecPath()
	if err != nil {
		log.Fatalln(err)
		c.Errors = append(c.Errors, err)
	}
	m.ExecPath = systems.ReplaceSplit(m.ExecPath)
	m.OSType = systems.GetGOOS()
	if m.OSType == "linux" {
		m.IsLinux = true
	} else {
		m.IsWindows = true
	}
	if m.MaxCPUThread > 0 {
		runtime.GOMAXPROCS(m.MaxCPUThread)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	}
	var routineFunc = func() error {
		var (
			err   error
			now   = time.Now().Local()
			upath = systems.ReplaceSplit(fmt.Sprintf("%s%s", m.RootPath, m.UploadPath))
			dpath = systems.ReplaceSplit(fmt.Sprintf("%s%s", m.RootPath, m.DownloadPath))
			name  = ""
		)
		_, err = systems.MustOpen(name, upath)
		Assert(err)
		_, err = systems.MustOpen(name, dpath)
		Assert(err)
		_, err = systems.MustOpen(name, dpath+systems.GetSplit()+time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.Local).Format("20060102150405")+systems.GetSplit())
		Assert(err)
		_, err = systems.MustOpen(name, upath+systems.GetSplit()+time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.Local).Format("20060102150405")+systems.GetSplit())
		Assert(err)
		_, err = systems.MustOpen(name, dpath+systems.GetSplit()+time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.Local).Format("20060102150405")+systems.GetSplit()+m.ThumbNailPath)
		Assert(err)
		return err
	}
	go times.Routine(start, 10, 24, routineFunc)
	return nil
}

func (m *App) Save(fileName, filePath string) (*os.File, error) {
	var (
		path   = systems.ReplaceSplit(fmt.Sprintf("%s%s", m.RootPath, filePath))
		f, err = systems.MustOpen(fileName, path+systems.GetSplit())
	)
	return f, err
}

func (m *App) Open(fileName, filePath string) (*os.File, error) {
	return os.Open(systems.ReplaceSplit(fmt.Sprintf("%s%s", m.RootPath, filePath)) + systems.GetSplit() + fileName)
}

func (m *App) Remove(filePath string) (err error) {
	return os.Remove(filePath)
}

func (m *App) Upload(path, fileName string, file multipart.File) (string, error) {
	var (
		split    = systems.GetSplit()
		savePath = systems.ReplaceSplit(fmt.Sprintf("%s%s", m.RootPath, m.UploadPath))
	)
	path = systems.ReplaceSplit(path)
	if path[len(path)-len(split):] != split {
		path = path + split
	}
	if savePath[len(savePath)-len(split):] != split && path[0:len(split)] != split {
		savePath = savePath + split + path
	} else if savePath[len(savePath)-len(split):] == split && path[0:len(split)] == split {
		savePath = savePath + path[len(split):]
	}
	if err := systems.MustWrite(fileName, savePath, file); err != nil {
		return "", nil
	}
	return savePath + fileName, nil
}

func (m *App) Download(savePath string) ([]byte, error) {
	text, err := ioutil.ReadFile(savePath)
	if err != nil {
		return nil, err
	}
	return text, nil
}

func (m *App) Starter(c *Content) error {
	return nil
}

func (m *App) Router(s *Server) {

}
