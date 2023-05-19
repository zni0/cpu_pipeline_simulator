package vcpu

import (
	"fmt"

	"github.com/zni0/cpu_pipeline_simulator/constants"
	"github.com/zni0/cpu_pipeline_simulator/lg"
	"github.com/zni0/cpu_pipeline_simulator/utils"
)

func (Cpu *CPU) RunE() {

	Logger := lg.Logger{Stage: "E",
		CycleNo: Cpu.CycleNo}
	Logger.Info("Started Execute")
	defer Logger.Info("Completed Execute")

	input := 2  // DE Latch
	output := 3 // EM Latch
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

	if Cpu.CyclesLeft[inputLatch.OPCode] > 1 {
		//Stall
		Logger.Info("Ongoing Instruction, Need to stall!")
		Cpu.WriteEnableSignal[input] = false // Disable self's write enable so as to not loose the instruction
		Cpu.AdjustWriteEnableSignals(input)
		if Cpu.WriteEnableSignal[output] {
			outputLatch.ValidBit = false
		}
		return
	} else if Cpu.CyclesLeft[inputLatch.OPCode] == 1 {
		Logger.Info("Instruction Completed in this cycle!")
	}

	instruction := inputLatch.Instruction
	pc := inputLatch.ProgramCounter
	destinationRegister := inputLatch.DestinationRegister
	source1 := inputLatch.Source1
	source2 := inputLatch.Source2
	sourceReg1 := inputLatch.SourceReg1
	sourceReg2 := inputLatch.SourceReg2
	literal := inputLatch.Literal
	opCode := inputLatch.OPCode
	Logger.Info("Read all inputs")
	utils.DoneAndWait(Cpu.ReadLatchWG) // Wait till all stages read registers
	// input latch / CPU registers should not be refered after this point
	Logger.Info(fmt.Sprintf("Instruction: %s", instruction))

	// Set latency of each operation
	if Cpu.CyclesLeft[inputLatch.OPCode] == 0 {
		switch opCode {
		case constants.MUL:
			Cpu.CyclesLeft[opCode] = 2
		default:
			Cpu.CyclesLeft[opCode] = 1
		}
	}

	// Logic for exec Stage
	var aluOutPut int
	switch opCode {
	case constants.ADD:
		aluOutPut = source1 + source2
	case constants.ADDI:
		aluOutPut = source1 + literal
	case constants.SUB:
		aluOutPut = source1 - source2
	case constants.SUBI:
		aluOutPut = source1 - literal
	case constants.MUL:
		aluOutPut = source1 * source2
	case constants.MOVC:
		aluOutPut = literal + 0
	case constants.BEQ, constants.BNE:
		Logger.Info("Control flow instruction")
		if (opCode == constants.BEQ && (source1 == source2)) ||
			(opCode == constants.BNE && (source1 != source2)) {
			Logger.Info("Branch Taken!")
			Cpu.AdjustWriteEnableSignals(input)
			Cpu.PCWG.Wait()                   // Wait till the Fetch state updates the PC and then over write it to branch
			Cpu.ProgramCounter = pc + literal //This logic is not inside the if as we know the next 2 stages won't stall
			// TODO: Should we move the PC update inside if?
			Logger.Info(fmt.Sprintf("PC set to %d", Cpu.ProgramCounter))
			if Cpu.WriteEnableSignal[output] {
				Logger.Info("Writing into E-M latch")
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
				Logger.Info("M stage is stalled, skipping write into E-M latch")
			}
			return
		}
	default:
		Logger.Error("Invalid Instruction")
	}

	if Cpu.CyclesLeft[opCode] > 1 {
		Cpu.CyclesLeft[opCode] -= 1
		Logger.Info(fmt.Sprintf("Running instruction: %s, %d cycles left",
			instruction, Cpu.CyclesLeft[opCode]))
		// Stall
		Logger.Info("Need to stall!")
		Cpu.WriteEnableSignal[input] = false // Disable self's write enable so as to not loose the instruction
		Cpu.AdjustWriteEnableSignals(input)
		if Cpu.WriteEnableSignal[output] {
			outputLatch.ValidBit = false
		}
		return
	}

	Cpu.AdjustWriteEnableSignals(input)
	// Write to memoryAccessStage
	if Cpu.WriteEnableSignal[output] {
		Cpu.CyclesLeft[opCode] = 0
		Logger.Info("Writing into E-M latch")
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
		Logger.Info("M stage is stalled, skipping write into E-M latch")
	}
	return
}
