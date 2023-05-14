package lg

import (
	"fmt"
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

func GetLog(CycleNo int) {
	logPerCycle := make(map[int][]Log)
	for i := 0; i < len(logs); i++ {
		logPerCycle[logs[i].CycleNo] = append(logPerCycle[logs[i].CycleNo], logs[i])
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println("CycleNo:", CycleNo)
	logPerStage := make(map[string][]string)
	for _, v := range logPerCycle[CycleNo] {
		logPerStage[v.Stage] = append(logPerStage[v.Stage],
			fmt.Sprintf("%s: %s", v.Level, v.Log))
	}

	stages := []string{"F", "D", "E", "M", "W"}
	for _, k := range stages {
		v := logPerStage[k]
		fmt.Println("Stage:", k)
		for j := 0; j < len(v); j++ {
			fmt.Println(v[j])
		}
		fmt.Println("------------")
	}
	fmt.Println("===========")
}

func GetAllLogs() {
	logPerCycle := make(map[int][]Log)
	for i := 0; i < len(logs); i++ {
		logPerCycle[logs[i].CycleNo] = append(logPerCycle[logs[i].CycleNo], logs[i])
	}

	for i := 0; i <= 1+len(logPerCycle); i++ {
		fmt.Println("CycleNo:", i)
		logPerStage := make(map[string][]string)
		for _, v := range logPerCycle[i] {
			logPerStage[v.Stage] = append(logPerStage[v.Stage],
				fmt.Sprintf("%s: %s", v.Level, v.Log))
		}

		stages := []string{"F", "D", "E", "M", "W"}
		for _, k := range stages {
			v := logPerStage[k]
			fmt.Println("Stage:", k)
			for j := 0; j < len(v); j++ {
				fmt.Println(v[j])
			}
			fmt.Println("------------")
		}
		fmt.Println("===========")
	}
}
