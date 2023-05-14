package vcpu

import (
	"sync"

	"github.com/zni0/cpu_pipeline_simulator/utils"
)

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
	go utils.AsyncSyncRun(&WG, Cpu.RunF)
	go utils.AsyncSyncRun(&WG, Cpu.RunD)
	go utils.AsyncSyncRun(&WG, Cpu.RunE)
	go utils.AsyncSyncRun(&WG, Cpu.RunM)
	go utils.AsyncSyncRun(&WG, Cpu.RunW)
	WG.Wait()
	// fmt.Println("---------------")
	valid := false
	for i := 1; i < 5; i++ {
		valid = valid || Cpu.Stages[i].ValidBit
	}
	// mx = mx - 1
	return valid
}

func (Cpu *CPU) AdjustWriteEnableSignals(input int) {
	utils.DoneAndWait(Cpu.SelfWriteEnableWG) // Wait till all stages complete computation, and set it's own write enable
	enabled := true
	for i := input; i < 5; i++ {
		Cpu.WriteEnableMu.Lock()
		enabled = enabled && Cpu.WriteEnableSignal[i]
		Cpu.WriteEnableMu.Unlock()
	}
	if !enabled {
		for i := 0; i <= input; i++ {
			Cpu.WriteEnableMu.Lock()
			Cpu.WriteEnableSignal[i] = false
			Cpu.WriteEnableMu.Unlock()
		}
	}
	utils.DoneAndWait(Cpu.WriteEnableCompleteWG) // Wait till all stages send Write enable signals to previous stages
}
