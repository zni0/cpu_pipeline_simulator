package utils

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/zni0/cpu_pipeline_simulator/constants"
)

func DecodeInstruction(instruction string) (int, int, int, int) {
	//ADD 1 1 3
	word := strings.Fields(instruction)
	opCode := word[0]

	// opCode, operand1, operand2, operand3, err := decodeInstruction(decodeStage.Instruction)
	switch opCode {
	case "ADD":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		s2, _ := strconv.Atoi(word[3])
		return constants.ADD, des, s1, s2
	case "ADDI":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		s2, _ := strconv.Atoi(word[3])
		return constants.ADDI, des, s1, s2
	case "SUB":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		s2, _ := strconv.Atoi(word[3])
		return constants.SUB, des, s1, s2
	case "SUBI":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		s2, _ := strconv.Atoi(word[3])
		return constants.SUBI, des, s1, s2
	case "MUL":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		s2, _ := strconv.Atoi(word[3])
		return constants.MUL, des, s1, s2
	case "MOVC":
		des, _ := strconv.Atoi(word[1])
		s1, _ := strconv.Atoi(word[2])
		return constants.MOVC, des, s1, -1
	// case "JUMP":
	// 	location, _ := strconv.Atoi(word[1])
	// 	return JUMP, location, -1, -1
	case "BEQ":
		s1, _ := strconv.Atoi(word[1])
		s2, _ := strconv.Atoi(word[2])
		jump, _ := strconv.Atoi(word[3])
		return constants.BEQ, s1, s2, jump
	case "BNE":
		s1, _ := strconv.Atoi(word[1])
		s2, _ := strconv.Atoi(word[2])
		jump, _ := strconv.Atoi(word[3])
		return constants.BNE, s1, s2, jump
	case "NOOP":
		return constants.NOOP, -1, -1, -1
	default:
		return -1, -1, -1, -1
	}

}

func DoneAndWait(wg *sync.WaitGroup) {
	wg.Done()
	wg.Wait()
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

func AsyncSyncRun(wg *sync.WaitGroup, f func()) {
	defer wg.Done()
	f()
}
