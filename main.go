package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/robfig/cron"

	"github.com/zwh8800/golang-mirror/spider"
)

func main() {
	defer glog.Flush()

	cronTab := cron.New()
	cronTab.AddFunc("@daily", spider.Go)
	cronTab.Start()
	glog.Infoln("server started")

	spider.Go()

	handleSignal()
	spider.WaitFinish()
	glog.Infoln("gracefully shutdown")
}

func handleSignal() {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	glog.Infoln("signal received")
}
