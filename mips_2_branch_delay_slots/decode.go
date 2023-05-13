package main

import "fmt"

func (Cpu *CPU) RunD() {

	Logger := Logger{Stage: "D",
		CycleNo: Cpu.CycleNo}
	Logger.Info("Started Decode")
	defer Logger.Info("Completed Decode")
	input := 1  // FD latch
	output := 2 // DE latch
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
	registerFilesATM := CloneMyRegFile(Cpu.RegisterFile)
	Logger.Info("Read all inputs")
	DoneAndWait(Cpu.ReadLatchWG) // Wait till all stages read registers
	// input latch / CPU registers should not be refered after this point
	Logger.Info(fmt.Sprintf("Instruction: %s", instruction))

	if instruction == "HALT" {
		Logger.Info("HALT Encountered")
		Cpu.Halt = true // Read in Fetch stage, after AdjustWriteEnableSignals.
		Cpu.AdjustWriteEnableSignals(input)
		if Cpu.WriteEnableSignal[output] {
			outputLatch.ValidBit = false
		}
		return
	}

	var opCode, source1, source2, literal, destination, sourceReg1, sourceReg2 int
	opCode, operand1, operand2, operand3 := decodeInstruction(instruction)
	switch opCode {
	case ADD, SUB, MUL:
		destination = operand1
		source1 = registerFilesATM[operand2].Value
		source2 = registerFilesATM[operand3].Value
		sourceReg1 = operand2
		sourceReg2 = operand3
		literal = -1
	case ADDI, SUBI:
		destination = operand1
		source1 = registerFilesATM[operand2].Value
		source2 = -1
		sourceReg1 = operand2
		sourceReg2 = -1
		literal = operand3
	case MOVC:
		destination = operand1
		source1 = -1
		source2 = -1
		sourceReg1 = -1
		sourceReg2 = -1
		literal = operand2
	case BEQ, BNE:
		source1 = registerFilesATM[operand1].Value
		source2 = registerFilesATM[operand2].Value
		sourceReg1 = operand1
		sourceReg2 = operand2
		literal = operand3
	case NOOP:
		destination = -1
		source1 = -1
		source2 = -1
		sourceReg1 = -1
		sourceReg2 = -1
		literal = -1
	default:
		Logger.Error("Unknown instruction")
	}

	if CheckIfStalled(sourceReg1, sourceReg2, destination, registerFilesATM) {
		// Need to stall
		Logger.Info("Need to stall!")
		Cpu.WriteEnableSignal[input] = false // Disable self's write enable so as to not loose the instruction
		Cpu.AdjustWriteEnableSignals(input)
		if Cpu.WriteEnableSignal[output] {
			outputLatch.ValidBit = false
		}
		return
	}

	Cpu.AdjustWriteEnableSignals(input)
	if Cpu.WriteEnableSignal[output] {
		Logger.Info("Writing into D-E latch")
		if opCode == NOOP {
			Logger.Info(fmt.Sprintf("NOOP encountered, setting ValidBit of DE latch to false"))
			outputLatch.ValidBit = false
			return
		}
		outputLatch.ValidBit = true
		outputLatch.OPCode = opCode
		outputLatch.DestinationRegister = destination
		outputLatch.Source1 = source1
		outputLatch.Source2 = source2
		outputLatch.SourceReg1 = sourceReg1
		outputLatch.SourceReg2 = sourceReg2
		outputLatch.Literal = literal
		outputLatch.Instruction = instruction
		outputLatch.ProgramCounter = pc
		// Set Score board
		Logger.Info("Setting Score-Board")
		if sourceReg2 != -1 {
			Cpu.RegisterFile[sourceReg1].InUse = true
		}
		if sourceReg2 != -1 {
			Cpu.RegisterFile[sourceReg2].InUse = true
		}
		if destination != -1 {
			Cpu.RegisterFile[destination].InUse = true
		}
		Cpu.RegisterFile[destination].InUse = true
	} else {
		Logger.Info("E stage is stalled, skipping write into D-E latch")
	}
}

func CheckIfStalled(sourceReg1, sourceReg2, destination int, registerFilesATM []Register) bool {
	source1InUse := false
	if sourceReg1 != -1 {
		source1InUse = registerFilesATM[sourceReg1].InUse

	}
	source2InUse := false
	if sourceReg2 != -1 {
		source2InUse = registerFilesATM[sourceReg2].InUse

	}
	destinationInUse := false
	if destination != -1 {
		destinationInUse = registerFilesATM[destination].InUse

	}
	return (destinationInUse || source1InUse || source2InUse)
}
