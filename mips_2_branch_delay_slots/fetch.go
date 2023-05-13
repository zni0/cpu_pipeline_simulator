package main

import "fmt"

func (Cpu *CPU) RunF() {

	Logger := Logger{Stage: "F",
		CycleNo: Cpu.CycleNo}
	Logger.Info("Started Fetch")
	defer Logger.Info("Completed Fetch")
	input := 0
	output := 1
	outputLatch := Cpu.Stages[output]

	if Cpu.Halt {
		Logger.Info("CPU Halted")
		DoneAndWait(Cpu.ReadLatchWG)
		Cpu.AdjustWriteEnableSignals(input)
		if Cpu.WriteEnableSignal[output] {
			outputLatch.ValidBit = false
		}
		// Cpu.PCWG.Done() // Not needed as if CPU is halted, no further branch is waiting to write into pc
		return
	}

	pc := Cpu.ProgramCounter
	Logger.Info("Read all inputs")
	DoneAndWait(Cpu.ReadLatchWG) // Wait till all stages read registers
	// input latch / CPU registers should not be read after this point

	instruction := Memory[pc]
	pc = pc + 1

	Logger.Info(fmt.Sprintf("Instruction: %s", instruction))
	Cpu.AdjustWriteEnableSignals(input)

	if Cpu.Halt {
		Logger.Info(fmt.Sprintf("HALT encountered, skipping setting valid bit of D to false"))
		outputLatch.ValidBit = false
		Cpu.PCWG.Done()
		return
	}
	if Cpu.WriteEnableSignal[output] {
		Logger.Info(fmt.Sprintf("Setting PC in F stage, may be updated if some Branch is taken"))
		Cpu.ProgramCounter = pc
		Cpu.PCWG.Done()
		Logger.Info(fmt.Sprintf("Writing into F-D latch"))
		outputLatch.ValidBit = true
		outputLatch.Instruction = instruction
		outputLatch.ProgramCounter = pc
	} else {
		Logger.Info("D stage is stalled, skipping update of PC and write into F-D latch")
		Cpu.PCWG.Done()
	}
}
