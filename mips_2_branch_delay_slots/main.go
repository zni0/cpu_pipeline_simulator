package main

import (
	"fmt"
	"sync"

	"github.com/zni0/cpu_pipeline_simulator/lg"
	"github.com/zni0/cpu_pipeline_simulator/vcpu"
)

// TODO: Not entirly correct, the execution will stop if all slots have NOOPs!
// Eg: "NOOP", "HALT" will terminalte one cycle early!
var mx int

func runEntireCode(Cpu *vcpu.CPU) {
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
	lg.GetAllLogs()
}

func runCycleByCycle(Cpu *vcpu.CPU) {
	j := 0
	for true {
		Cpu.CycleNo = j
		valid := Cpu.RunCycle()
		var ins string
		fmt.Scanln(&ins)
		lg.GetLog(Cpu.CycleNo)
		j++
		if !valid {
			break
		}
	}
}

func Init() vcpu.CPU {
	var readLatch, writeEnable, selfWriteEnable, pc sync.WaitGroup
	writeEnableMu := sync.Mutex{}

	registerFile := [32]vcpu.Register{}
	inputSignals := make([]bool, 5)
	cyclesLeft := make([]int, 100)
	fetchStage := vcpu.Stage{
		ValidBit: true,
	}
	decodeStage := vcpu.Stage{}
	execStage := vcpu.Stage{}
	memoryStage := vcpu.Stage{}
	writeBackStage := vcpu.Stage{}
	stages := [...]*vcpu.Stage{&fetchStage, &decodeStage, &execStage, &memoryStage, &writeBackStage}
	Cpu := vcpu.CPU{
		CycleNo:               0,
		ProgramCounter:        0,
		RegisterFile:          registerFile[:],
		Stages:                stages[:],
		ReadLatchWG:           &readLatch,
		WriteEnableCompleteWG: &writeEnable,
		SelfWriteEnableWG:     &selfWriteEnable,
		PCWG:                  &pc,
		Halt:                  false,
		WriteEnableSignal:     inputSignals,
		WriteEnableMu:         &writeEnableMu,
		CyclesLeft:            cyclesLeft,
	}
	return Cpu
}

func main() {

	Cpu := Init()
	vcpu.LoadCode()
	// runEntireCode(&Cpu)
	runCycleByCycle(&Cpu)
	fmt.Println(Cpu.RegisterFile)
}
