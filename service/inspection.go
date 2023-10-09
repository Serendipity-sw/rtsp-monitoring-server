package service

import (
	"fmt"
	"github.com/Serendipity-sw/gutil"
	"github.com/guotie/deferinit"
	"github.com/swgloomy/gutil/glog"
	"os"
	"path/filepath"
	"rtsp-monitoring-server/store"
	"strings"
	"sync"
	"time"
)

func Init() {
	deferinit.AddRoutine(inspection)
}

func inspection(ch chan struct{}, wg *sync.WaitGroup) {
	var (
		jsTmr     *time.Timer
		minNumber = time.Duration(5)
		fileName  string
		extension string
	)
	go func() {
		<-ch
		jsTmr.Stop()
		wg.Done()
	}()
	jsTmr = time.NewTimer(minNumber * time.Minute)
	<-jsTmr.C
	for {
		glog.Info("inspection start! \n")
		timeNow := time.Now()
		var storeVideoList []string
		store.MonitoringSync.Lock()
		for _, monitor := range store.MonitoringList {
			duration := timeNow.Sub(monitor.StartTime)
			if duration.Minutes() > 5 {
				monitor.ExitSignal <- struct{}{}
			}
			fileName = filepath.Base(monitor.FileName)
			extension = filepath.Ext(fileName)
			storeVideoList = append(storeVideoList, strings.TrimSuffix(fileName, extension))
		}
		store.MonitoringSync.Unlock()
		removeWidthOutFileName(storeVideoList)
		glog.Info("inspection start run success! \n")
		jsTmr.Reset(minNumber * time.Minute)
		glog.Info("inspection is waiting! \n")
		<-jsTmr.C
	}
}

func removeWidthOutFileName(fileNameList []string) {
	fileNameListIn, err := gutil.GetMyAllFileByDir(store.Asset)
	if err != nil {
		glog.Error("service removeWidthOutFileName GetMyAllFileByDir run err! fileNameList: %+v err: %+v \n", fileNameList, err)
		return
	}
	var (
		removeFilePath string
	)
	for _, fileName := range *fileNameListIn {
		for _, storeFileItem := range fileNameList {
			if strings.Index(fileName, storeFileItem) == -1 {
				removeFilePath = fmt.Sprintf("%s/%s", store.Asset, fileName)
				err = os.Remove(removeFilePath)
				if err != nil {
					glog.Error("service removeWidthOutFileName remove file err! removeFilePath: %s err: %+v \n", removeFilePath, err)
				} else {
					glog.Info("service removeWidthOutFileName remove file success! removeFilePath: %s \n", removeFilePath)
				}
			}
		}
	}
	glog.Info("service removeWidthOutFileName run success! \n")
}

func SetMonitorStartTime(channel string) {
	store.MonitoringSync.Lock()
	defer store.MonitoringSync.Unlock()
	item, ok := store.MonitoringList[channel]
	if ok {
		item.StartTime = time.Now()
		store.MonitoringList[channel] = item
	}
}
