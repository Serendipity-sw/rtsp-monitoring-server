package config

import (
	"fmt"
	"github.com/guotie/config"
	"os"
	"rtsp-monitoring-server/store"
)

func Init() {
	store.ListenPort = config.GetStringDefault("listenPort", ":8080")
	store.RootPrefix = config.GetStringDefault("rootPrefix", "")
	store.RtspFormat = config.GetStringDefault("rtspFormat", "")
	store.LogsDir = config.GetStringDefault("logsDir", "./logs")
	_, err := os.Stat(store.LogsDir)
	dirExisted := err == nil || os.IsExist(err)
	if !dirExisted {
		err = os.Mkdir(store.LogsDir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir run err! logsDir: %s err: %+v \n", store.LogsDir, err)
		}
	}
	store.Template = config.GetStringDefault("template", "./template")
	_, err = os.Stat(store.Template)
	dirExisted = err == nil || os.IsExist(err)
	if !dirExisted {
		err = os.Mkdir(store.Template, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir run err! template: %s err: %+v \n", store.Template, err)
		}
	}
	store.Asset = config.GetStringDefault("asset", "./asset")
	_, err = os.Stat(store.Asset)
	dirExisted = err == nil || os.IsExist(err)
	if !dirExisted {
		err = os.Mkdir(store.Asset, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir run err! asset: %s err: %+v \n", store.Asset, err)
		}
	}
}
