package spider

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"runtime/debug"
	"sync"

	"github.com/golang/glog"
	"gopkg.in/go-playground/pool.v1"

	"github.com/zwh8800/golang-mirror/conf"
	"github.com/zwh8800/golang-mirror/model"
	"github.com/zwh8800/golang-mirror/workspace"
)

var waitFinish sync.WaitGroup

type golangIndex struct {
	FileList []*model.File `xml:"Contents"`
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
	p := pool.NewPool(conf.Conf.Golang.MaxDownloadThread, len(index.FileList))

	for _, file := range index.FileList {
		if file.LastModified.After(conf.Conf.Golang.Earliest) {
			p.Queue(func(job *pool.Job) {
				file := job.Params()[0].(*model.File)
				if workspace.IsFileDirty(file) {
					glog.Infoln("file", file.Key, "dirty, updating")
					if err := updateFile(file); err != nil {
						glog.Errorln("error when updating", file.Key, ":", err)
					}
				} else {
					glog.Infoln("file", file.Key, "not modified")
				}
			}, file)
		}
	}
	for result := range p.Results() {
		err, ok := result.(*pool.ErrRecovery)
		if ok {
			glog.Errorln(err)
			return err
		}
	}

	return nil
}

func Go() {
	waitFinish.Add(1)
	defer waitFinish.Done()
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

func WaitFinish() {
	waitFinish.Wait()
}
