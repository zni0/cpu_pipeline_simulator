package main

import "fmt"

func (Cpu *CPU) RunE() {

	Logger := Logger{Stage: "E",
		CycleNo: Cpu.CycleNo}
	Logger.Info("Started Execute")
	defer Logger.Info("Completed Execute")
	input := 2  // DE Latch
	output := 3 // EM Latch
	inputLatch := Cpu.Stages[input]
	outputLatch := Cpu.Stages[output]

	if inputLatch.ValidBit == false {
		Logger.Info("Invalid Instruction (NOOP)")
		DoneAndWait(Cpu.ReadLatchWG) // Wait till all stages read registers
		Cpu.AdjustWriteEnableSignals(input)
		if Cpu.WriteEnableSignal[output] {
			outputLatch.ValidBit = false
		}
		return
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
	DoneAndWait(Cpu.ReadLatchWG) // Wait till all stages read registers
	// input latch / CPU registers should not be refered after this point
	Logger.Info(fmt.Sprintf("Instruction: %s", instruction))

	// Logic for exec Stage
	var aluOutPut int
	switch opCode {
	case ADD:
		aluOutPut = source1 + source2
	case ADDI:
		aluOutPut = source1 + literal
	case SUB:
		aluOutPut = source1 - source2
	case SUBI:
		aluOutPut = source1 - literal
	case MUL:
		aluOutPut = source1 * source2
	case MOVC:
		aluOutPut = literal + 0
	case BEQ, BNE:
		Logger.Info("Control flow instruction")
		if (opCode == BEQ && (source1 == source2)) ||
			(opCode == BNE && (source1 != source2)) {
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

	// No need to stall so WriteEnabled stays true!
	// In case of multiple cycle operations, we may need to stall and set WriteEnableSignal of input to 0 and stall
	Cpu.AdjustWriteEnableSignals(input)
	// Write to memoryAccessStage
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
