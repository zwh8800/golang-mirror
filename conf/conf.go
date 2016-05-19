package conf

import (
	"flag"
	"path/filepath"
	"time"

	"github.com/golang/glog"
	"gopkg.in/gcfg.v1"
)

type config struct {
	Golang struct {
		IndexPage string
		Earliest  time.Time
	}
	WeixinPush struct {
		Url string
	}
	WorkSpace struct {
		Path string
	}
}

var Conf config

func ReadConf(filename string) error {
	absFile := filename
	if !filepath.IsAbs(filename) {
		var err error
		absFile, err = filepath.Abs(filename)
		if err != nil {
			return err
		}
	}
	return gcfg.ReadFileInto(&Conf, absFile)
}

func init() {
	configFilename := flag.String("config", "golang-mirror.gcfg", "specify a config file")
	flag.Parse()
	glog.Infoln("configuring...")

	if err := ReadConf(*configFilename); err != nil {
		glog.Fatalf("error occored: %s", err)
		panic(err)
	}
}
