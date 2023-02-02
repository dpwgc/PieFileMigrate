package base

import (
	"sync"
	"time"
)

var WorkerMonitorMap sync.Map

type WorkerMonitorModel struct {
	LastMigrateStartTime time.Time `json:"lastMigrateStartTime"`
	LastMigrateEndTime   time.Time `json:"lastMigrateEndTime"`
}
