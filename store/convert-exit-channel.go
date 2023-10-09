package store

import (
	class "rtsp-monitoring-server/struct"
	"sync"
)

var (
	MonitoringList = make(map[string]class.Monitor)
	MonitoringSync sync.RWMutex
)
