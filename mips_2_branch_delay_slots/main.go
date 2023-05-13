package main

import (
	"fmt"
	"sync"
)

// TODO: Not entirly correct, the execution will stop if all slots have NOOPs!
// Eg: "NOOP", "HALT" will terminalte one cycle early!
var mx int

func (Cpu *CPU) RunCycle() bool {

	Cpu.ReadLatchWG.Add(5)
	Cpu.SelfWriteEnableWG.Add(5)
	Cpu.WriteEnableCompleteWG.Add(5)
	Cpu.PCWG.Add(1)
	for i := 0; i < 5; i++ {
		Cpu.WriteEnableSignal[i] = true
	}
	var WG sync.WaitGroup
	WG.Add(5)
	go AsyncSyncRun(&WG, Cpu.RunF)
	go AsyncSyncRun(&WG, Cpu.RunD)
	go AsyncSyncRun(&WG, Cpu.RunE)
	go AsyncSyncRun(&WG, Cpu.RunM)
	go AsyncSyncRun(&WG, Cpu.RunW)
	WG.Wait()
	// fmt.Println("---------------")
	valid := false
	for i := 1; i < 5; i++ {
		valid = valid || Cpu.Stages[i].ValidBit
	}
	mx = mx - 1
	return valid
}

func runEntireCode(Cpu *CPU) {
	i := 0
	mx = 20000
	for mx > 0 {
		Cpu.CycleNo = i
		running := Cpu.RunCycle()
		if !running {
			break
		}
		i++
	}
	fmt.Println(Cpu.RegisterFile)
	fmt.Println(Cpu.CycleNo)
	// fmt.Println(Cpu.RegisterFile[2])
	// fmt.Println(Cpu.RegisterFile[3])

	logPerCycle := make(map[int][]Log)
	for i := 0; i < len(logs); i++ {
		logPerCycle[logs[i].CycleNo] = append(logPerCycle[logs[i].CycleNo], logs[i])
	}

	// fmt.Println(logPerCycle[0])
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

func runCycleByCycle(Cpu *CPU) {
	j := 0
	for true {
		Cpu.CycleNo = j
		valid := Cpu.RunCycle()
		var ins string
		fmt.Scanln(&ins)
		logPerCycle := make(map[int][]Log)
		for i := 0; i < len(logs); i++ {
			logPerCycle[logs[i].CycleNo] = append(logPerCycle[logs[i].CycleNo], logs[i])
		}

		// fmt.Println(logPerCycle[0])
		fmt.Println("CycleNo:", j)
		logPerStage := make(map[string][]string)
		for _, v := range logPerCycle[j] {
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
		j++
		if !valid {
			break
		}
	}
}
func main() {

	Cpu := CPU_init()
	loadCode()
	// runEntireCode(&Cpu)
	runCycleByCycle(&Cpu)
}
