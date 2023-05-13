package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
)

var writeEnableMu = &sync.Mutex{}

func decodeInstruction(instruction string) (int, int, int, int) {
	//ADD 1 1 3
	word := strings.Fields(instruction)
	opCode := word[0]

	// opCode, operand1, operand2, operand3, err := decodeInstruction(decodeStage.Instruction)
	switch opCode {
	case "ADD":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		s2, _ := strconv.Atoi(word[3])
		return ADD, des, s1, s2
	case "ADDI":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		s2, _ := strconv.Atoi(word[3])
		return ADDI, des, s1, s2
	case "SUB":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		s2, _ := strconv.Atoi(word[3])
		return SUB, des, s1, s2
	case "SUBI":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		s2, _ := strconv.Atoi(word[3])
		return SUBI, des, s1, s2
	case "MUL":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		s2, _ := strconv.Atoi(word[3])
		return MUL, des, s1, s2
	case "MOVC":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		return MOVC, des, s1, -1
	// case "JUMP":
	// 	location, _ := strconv.Atoi(word[1])
	// 	return JUMP, location, -1, -1
	case "BEQ":
		s1, _ := strconv.Atoi(word[1])
		s2, _ := strconv.Atoi(word[2])
		jump, _ := strconv.Atoi(word[3])
		return BEQ, s1, s2, jump
	case "BNE":
		s1, _ := strconv.Atoi(word[1])
		s2, _ := strconv.Atoi(word[2])
		jump, _ := strconv.Atoi(word[3])
		return BNE, s1, s2, jump
	case "NOOP":
		return NOOP, -1, -1, -1
	default:
		return -1, -1, -1, -1
	}

}

func DoneAndWait(wg *sync.WaitGroup) {
	wg.Done()
	wg.Wait()
}

func (Cpu *CPU) AdjustWriteEnableSignals(input int) {
	DoneAndWait(Cpu.SelfWriteEnableWG) // Wait till all stages complete computation, and set it's own write enable
	enabled := true
	for i := input; i < 5; i++ {
		writeEnableMu.Lock()
		enabled = enabled && Cpu.WriteEnableSignal[i]
		writeEnableMu.Unlock()
	}
	if !enabled {
		for i := 0; i <= input; i++ {
			writeEnableMu.Lock()
			Cpu.WriteEnableSignal[i] = false
			writeEnableMu.Unlock()
		}
	}
	DoneAndWait(Cpu.WriteEnableCompleteWG) // Wait till all stages send Write enable signals to previous stages
}

func CloneMyRegFile(orig []Register) []Register {
	origJSON, _ := json.Marshal(orig)
	clone := []Register{}
	_ = json.Unmarshal(origJSON, &clone)
	return clone
}

func AsyncSyncRun(wg *sync.WaitGroup, f func()) {
	defer wg.Done()
	f()
}

func CPU_init() CPU {
	var readLatch, writeEnable, selfWriteEnable, pc sync.WaitGroup
	registerFile := [32]Register{}
	inputSignals := make([]bool, 5)
	fetchStage := Stage{
		ValidBit:    true,
		WriteEnable: true,
	}
	decodeStage := Stage{WriteEnable: true}
	execStage := Stage{WriteEnable: true}
	memoryStage := Stage{WriteEnable: true}
	writeBackStage := Stage{WriteEnable: true}
	stages := [...]*Stage{&fetchStage, &decodeStage, &execStage, &memoryStage, &writeBackStage}
	Cpu := CPU{
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
	}
	return Cpu

}

func InterfaceToString(data interface{}) string {
	switch data.(type) {
	case string:
		return data.(string)
	default:
		dataBytes, err := json.Marshal(data)
		if err != nil {
			dataBytes = []byte{}
		}
		return string(dataBytes)
	}
}
