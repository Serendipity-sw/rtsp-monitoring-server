package main

import (
	"flag"
	"fmt"
	"github.com/Serendipity-sw/gutil"
	readConfig "github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/swgloomy/gutil/glog"
	"os"
	"os/signal"
	"rtsp-monitoring-server/config"
	"rtsp-monitoring-server/router"
	"rtsp-monitoring-server/service"
	"rtsp-monitoring-server/store"
	"syscall"
)

var (
	pidStrPath = "./marketing-procurement-system.pid"
	debugFlag  = flag.Bool("d", false, "debug mode")
	configFn   = flag.String("config", "./config.json", "config file path")
)

func main() {
	flag.Parse()

	err := readConfig.ReadCfg(*configFn)
	if err != nil {
		fmt.Printf("main ReadCfg read err! filePath: %s err: %+v \n", *configFn, err.Error())
		return
	}

	config.Init()

	serverRun(*debugFlag)

	c := make(chan os.Signal, 1)
	gutil.WritePid(pidStrPath)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	//信号等待
	<-c
	fmt.Println("main exit application!")
	serverExit()
}

func serverRun(debug bool) {
	gutil.LogInit(debug, store.LogsDir)

	store.Init()
	fmt.Println("store init successfully!")

	service.Init()
	fmt.Println("service init successfully!")

	gutil.SetCPUUseNumber(0)
	fmt.Println("set many cpu successfully!")

	deferinit.InitAll()
	fmt.Println("init all module successfully!")

	deferinit.RunRoutines()
	fmt.Println("init all run successfully!")

	router.Init(debug)
	fmt.Println("ginInit run successfully!")
}

func serverExit() {
	deferinit.StopRoutines()
	fmt.Println("stop routine successfully!")

	deferinit.FiniAll()
	fmt.Println("stop all modules successfully!")

	glog.Close()

	os.Exit(0)
}
