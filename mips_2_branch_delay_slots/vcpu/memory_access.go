package vcpu

import (
	"fmt"

	"github.com/zni0/cpu_pipeline_simulator/lg"
	"github.com/zni0/cpu_pipeline_simulator/utils"
)

func (Cpu *CPU) RunM() {

	Logger := lg.Logger{Stage: "M",
		CycleNo: Cpu.CycleNo}
	Logger.Info("Started MemoryAccess")
	defer Logger.Info("Completed MemoryAccess")
	input := 3  // EM Latch
	output := 4 // MW Latch
	inputLatch := Cpu.Stages[input]
	outputLatch := Cpu.Stages[output]

	if inputLatch.ValidBit == false {
		Logger.Info("Invalid Instruction (NOOP)")
		utils.DoneAndWait(Cpu.ReadLatchWG) // Wait till all stages read registers
		Cpu.AdjustWriteEnableSignals(input)
		if Cpu.WriteEnableSignal[output] {
			outputLatch.ValidBit = false
		}
		return
	}

	instruction := inputLatch.Instruction
	_ = inputLatch.ProgramCounter
	destinationRegister := inputLatch.DestinationRegister
	source1 := inputLatch.Source1
	source2 := inputLatch.Source2
	sourceReg1 := inputLatch.SourceReg1
	sourceReg2 := inputLatch.SourceReg2
	literal := inputLatch.Literal
	aluOutPut := inputLatch.ALUOutPut
	Logger.Info("Read all inputs")
	utils.DoneAndWait(Cpu.ReadLatchWG) // Wait till all stages read registers
	// input latch / CPU registers should not be refered after this point
	Logger.Info(fmt.Sprintf("Instruction: %s", instruction))

	// Logic for exec Stage
	// TODO

	// No need to stall so WriteEnabled stays true!
	// In case of multiple cycle operations, we may need to stall and set WriteEnableSignal of input to 0 and stall
	Cpu.AdjustWriteEnableSignals(input)
	// Write to memoryAccessStage
	if Cpu.WriteEnableSignal[output] {
		Logger.Info("Writing into M-W latch")
		outputLatch.ValidBit = true
		outputLatch.ALUOutPut = aluOutPut
		outputLatch.DestinationRegister = destinationRegister
		outputLatch.Source1 = source1
		outputLatch.Source2 = source2
		outputLatch.SourceReg1 = sourceReg1
		outputLatch.SourceReg2 = sourceReg2
		outputLatch.Literal = literal
		outputLatch.Instruction = instruction
	} else {
		Logger.Info("W stage is stalled, skipping write into M-W latch")
	}

	return
}
