package main

import "fmt"

func (Cpu *CPU) RunW() {

	Logger := Logger{Stage: "W",
		CycleNo: Cpu.CycleNo}
	input := 4
	inputLatch := Cpu.Stages[input]

	if inputLatch.ValidBit == false {
		Logger.Info("Invalid Instruction (NOOP)")
		DoneAndWait(Cpu.ReadLatchWG) // Wait till all stages read registers
		Cpu.AdjustWriteEnableSignals(input)
		return
	}

	// instruction := inputLatch.Instruction
	destination := inputLatch.DestinationRegister
	sourceReg1 := inputLatch.SourceReg1
	sourceReg2 := inputLatch.SourceReg2
	aluOutPut := inputLatch.ALUOutPut
	instruction := inputLatch.Instruction
	Logger.Info("Read all inputs")
	DoneAndWait(Cpu.ReadLatchWG) // Wait till all stages read registers
	// input latch / CPU registers should not be refered after this point
	Logger.Info(fmt.Sprintf("Instruction: %s", instruction))

	// No logic in Write back stage

	// No need to stall so WriteEnabled stays true!
	// In case of multiple cycle operations, we may need to stall and set WriteEnableSignal of input to 0 and stall
	Cpu.AdjustWriteEnableSignals(input)

	// Logic for writeBackStage
	Cpu.RegisterFile[destination].Value = aluOutPut
	// Set Score board
	if sourceReg1 != -1 {
		Cpu.RegisterFile[sourceReg1].InUse = false
	}
	if sourceReg2 != -1 {
		Cpu.RegisterFile[sourceReg2].InUse = false
	}
	Cpu.RegisterFile[destination].InUse = false
}
