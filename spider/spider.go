package spider

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"runtime/debug"

	"github.com/golang/glog"

	"github.com/zwh8800/golang-mirror/conf"
	"github.com/zwh8800/golang-mirror/model"
	"github.com/zwh8800/golang-mirror/workspace"
)

type golangIndex struct {
	FileList []model.File `xml:"Contents"`
}

func updateFile(file *model.File) error {
	u, err := url.Parse(conf.Conf.Golang.IndexPage)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, file.Key)

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := workspace.InsertOrUpdateFile(file, resp.Body); err != nil {
		return err
	}
	return nil
}

func spiderIndex() error {
	resp, err := http.Get(conf.Conf.Golang.IndexPage)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	index := golangIndex{}
	if err := xml.Unmarshal(data, &index); err != nil {
		return err
	}
	for _, file := range index.FileList {
		go func() {
			if workspace.IsFileDirty(&file) {
				glog.Infoln("file", file.Key, "dirty, updating")
				if err := updateFile(&file); err != nil {
					glog.Errorln(err)
				}
			} else {
				glog.Infoln("file", file.Key, "not modified")
			}
		}()
	}

	return nil
}

func Go() {
	glog.Info("golang spider started")
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln("panic in spider recovered:", err, string(debug.Stack()))
		}
	}()

	if err := spiderIndex(); err != nil {
		glog.Errorln(err)
	}

	glog.Info("golang spider finished")
}
