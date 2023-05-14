package vcpu

import "encoding/json"

func CloneMyRegFile(orig []Register) []Register {
	origJSON, _ := json.Marshal(orig)
	clone := []Register{}
	_ = json.Unmarshal(origJSON, &clone)
	return clone
}
