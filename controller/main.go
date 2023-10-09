package controller

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/swgloomy/gutil/glog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"rtsp-monitoring-server/store"
	class "rtsp-monitoring-server/struct"
	"strings"
	"time"
)

func Open(c *gin.Context) {
	monitoringChannel := strings.TrimSpace(c.Query("monitoringChannel"))
	if monitoringChannel == "" {
		glog.Error("controller Open monitoringChannel is empty! \n")
		c.JSON(http.StatusOK, class.ResultStruct{Code: 500, Msg: "监控通道不能为空"})
		return
	}
	var (
		ch = make(chan string)
	)
	go monitoring(monitoringChannel, ch)
	fileName := <-ch
	glog.Info("controller Open run success! monitoringChannel: %s \n", monitoringChannel)
	c.JSON(http.StatusOK, class.ResultStruct{Code: 200, Data: fmt.Sprintf("%s?channel=%s", fileName, monitoringChannel), Msg: "监控开启"})
}

func monitoring(monitoringChannel string, ch chan string) {
	store.MonitoringSync.RLock()
	monitorItem, ok := store.MonitoringList[monitoringChannel]
	store.MonitoringSync.RUnlock()
	if ok {
		ch <- monitorItem.FileName
		return
	}
	uuidStr := uuid.NewString()
	mainFile := fmt.Sprintf("%s.m3u8", uuidStr)
	filePathStr := fmt.Sprintf("%s/%s/%s", store.ExeDir, store.Asset, mainFile)
	cmd := exec.Command("ffmpeg", "-rtsp_transport", "tcp", "-hwaccel", "cuda", "-i", fmt.Sprintf(store.RtspFormat, monitoringChannel), "-segment_atclocktime", "1", "-c:v", "copy", "-s", "720*576", "-r", "50", "-crf", "28", "-tune", "zerolatency", "-vcodec", "libx264", "-b:v", "400k", "-threads", "10", "-preset", "ultrafast", "-y", filePathStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		glog.Error("controller monitoring Command run err! monitoringChannel: %s err: %+v \n", monitoringChannel, err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		glog.Error("controller monitoring StderrPipe run err! monitoringChannel: %s err: %+v \n", monitoringChannel, err)
		return
	}
	if err = cmd.Start(); err != nil {
		glog.Error("controller monitoring cmd Start err! monitoringChannel: %s err: %+v \n", monitoringChannel, err)
		return
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			glog.Info("controller monitoring NewScanner run success! stdout: %s \n", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		findExit := false
		for scanner.Scan() {
			textStr := scanner.Text()
			if !findExit && strings.Index(textStr, "Opening") > -1 && strings.Index(textStr, store.ExeDir) > -1 && strings.Index(textStr, ".ts'") > -1 {
				findExit = true
				ch <- mainFile
			}
			glog.Info("controller monitoring NewScanner run success! stderr: %s \n", textStr)
		}
	}()
	var (
		stopChan = make(chan struct{})
	)
	store.MonitoringSync.Lock()
	store.MonitoringList[monitoringChannel] = class.Monitor{
		FileName:   mainFile,
		StartTime:  time.Now(),
		ExitSignal: stopChan,
	}
	store.MonitoringSync.Unlock()
	<-stopChan
	store.MonitoringSync.Lock()
	delete(store.MonitoringList, monitoringChannel)
	store.MonitoringSync.Unlock()
	taskKillCmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", cmd.Process.Pid))
	err = taskKillCmd.Run()
	if err != nil {
		glog.Error("controller monitoring taskKillCmd run err! err: %+v \n", err)
	} else {
		glog.Info("controller monitoring taskKillCmd run success! \n")
	}
	pattern := fmt.Sprintf("%s/%s*", store.Asset, uuidStr)
	files, err := filepath.Glob(pattern)
	if err != nil {
		glog.Error("controller monitoring filepath Glob run err! monitoringChannel: %s pattern: %s err: %+v \n", monitoringChannel, pattern, err)
	} else {
		for _, file := range files {
			err = os.Remove(file)
			if err != nil {
				glog.Error("controller monitoring file remove err! monitoringChannel: %s file: %s pattern: %s err: %+v \n", monitoringChannel, file, pattern, err)
			}
		}
		glog.Info("controller monitoring file remove run success! files: %+v monitoringChannel: %s pattern: %s  \n", files, monitoringChannel, pattern)
	}
	glog.Info("controller monitoring run success! monitoringChannel: %s \n", monitoringChannel)
}

func Close(c *gin.Context) {
	monitoringChannel := strings.TrimSpace(c.Query("monitoringChannel"))
	if monitoringChannel == "" {
		glog.Error("controller Close monitoringChannel is empty! \n")
		c.JSON(http.StatusOK, class.ResultStruct{Code: 500, Msg: "监控通道不能为空"})
		return
	}
	store.MonitoringSync.RLock()
	defer store.MonitoringSync.RUnlock()
	item, ok := store.MonitoringList[monitoringChannel]
	if ok {
		item.ExitSignal <- struct{}{}
	}
	glog.Info("controller Close run success! monitoringChannel: %s \n", monitoringChannel)
	c.JSON(http.StatusOK, class.ResultStruct{Code: 200, Msg: "关闭成功"})
}
