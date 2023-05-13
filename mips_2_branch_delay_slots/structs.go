package main

import (
	"sync"
)

var Memory []string

const (
	ADD int = iota
	ADDI
	SUB
	SUBI
	MUL
	MOVC
	BEQ
	BNE
	JUMP
	NOOP
)

type CPU struct {
	CycleNo               int // Clock cycles elasped
	ProgramCounter        int // Current program counter
	RegisterFile          []Register
	Stages                []*Stage
	WriteEnableSignal     []bool
	ReadLatchWG           *sync.WaitGroup
	SelfWriteEnableWG     *sync.WaitGroup
	WriteEnableCompleteWG *sync.WaitGroup
	PCWG                  *sync.WaitGroup
	Halt                  bool
}

type Register struct {
	Value int
	InUse bool
}

type Stage struct {
	Instruction        string // Instruction
	CyclesRemaining    int    // For instructions needing more than one cycle
	OnGoingInstruction bool

	ProgramCounter      int        // Program Counter
	Source1             int        // Source-1
	Source2             int        // Source-2
	SourceReg1          int        // Source-1
	SourceReg2          int        // Source-2
	DestinationRegister int        // Destination Register Address
	Literal             int        // Literal Value
	OPCode              int        // Operation Code
	ValidBit            bool       // Bit to indecate NoOp instruction
	ALUOutPut           int        // ALU's output in exec stage
	LatchWriteEnable    *chan bool // NextStage will send in if latch needs to be retained (for stall conditions)
	WriteEnable         bool
	NoOp                bool
}
