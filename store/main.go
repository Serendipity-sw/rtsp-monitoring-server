package store

import (
	"github.com/swgloomy/gutil/glog"
	"os"
	"path/filepath"
	"strings"
)

func Init() {
	exePath, err := os.Executable()
	if err != nil {
		glog.Error("store Init Executable run err! err: %+v \n", err)
		return
	}
	ExeDir = filepath.Dir(exePath)
	ExeDir = strings.ReplaceAll(ExeDir, "\\", "/")
	glog.Info("store Init run success! exeDir: %s \n", ExeDir)
}
