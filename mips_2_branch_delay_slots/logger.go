package main

import (
	"sync"
	"time"
)

var logs []Log
var logMu = &sync.Mutex{}

type Logger struct {
	Identifier string
	LogHistory string
	Stage      string
	CycleNo    int
}

type Log struct {
	CycleNo int
	Stage   string
	Time    time.Time
	Log     string
	Level   string
}

func (l Logger) Info(log string) {
	logObj := Log{
		CycleNo: l.CycleNo,
		Stage:   l.Stage,
		Time:    time.Now(),
		Log:     log,
		Level:   "Info",
	}
	logMu.Lock()
	logs = append(logs, logObj)
	logMu.Unlock()
}
func (l Logger) Error(log string) {
	logObj := Log{
		CycleNo: l.CycleNo,
		Stage:   l.Stage,
		Time:    time.Now(),
		Log:     log,
		Level:   "Error",
	}
	logMu.Lock()
	logs = append(logs, logObj)
	logMu.Unlock()
}
func (l Logger) Debug(log string) {
	logObj := Log{
		CycleNo: l.CycleNo,
		Stage:   l.Stage,
		Time:    time.Now(),
		Log:     log,
		Level:   "Debug",
	}
	logMu.Lock()
	logs = append(logs, logObj)
	logMu.Unlock()
}

// {
// 	"CycleNo": 1,
// 	"Stage": "E",
// 	"Time": 00102,
// 	"log": "Computing",
// 	"Level": "D"
// }
