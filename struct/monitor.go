package class

import "time"

type Monitor struct {
	FileName   string `json:"file_name"`
	ExitSignal chan struct{}
	StartTime  time.Time
}
