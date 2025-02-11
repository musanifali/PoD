package blockchain

import (
	"fmt"
)

// PoDConsensus represents the Proof-of-Data consensus mechanism
type PoDConsensus struct {
	Validators []*Validator
}

// NewPoDConsensus initializes PoD with validator nodes
func NewPoDConsensus() *PoDConsensus {
	return &PoDConsensus{
		Validators: []*Validator{},
	}
}

// RegisterValidator registers a validator node (Prevents duplicates)
func (pod *PoDConsensus) RegisterValidator(validator *Validator) {
	// ✅ Prevent duplicate validators
	for _, v := range pod.Validators {
		if v.ID == validator.ID {
			fmt.Println("⚠ Validator already registered:", validator.ID)
			return
		}
	}

	pod.Validators = append(pod.Validators, validator)
	fmt.Println("✅ Validator registered:", validator.ID)
}

// ValidateBlock ensures a block is approved by validators before adding
func (pod *PoDConsensus) ValidateBlock(block Block, blockchain *Blockchain) bool {
	if len(pod.Validators) == 0 {
		fmt.Println("❌ No validators registered! Block cannot be approved.")
		return false
	}

	approvedVotes := 0
	requiredVotes := (len(pod.Validators) * 75) / 100 // ✅ Requires 75% approval

	// Validators approve the block
	for _, validator := range pod.Validators {
		if validator.ApproveBlock(block) {
			approvedVotes++
		}
	}

	fmt.Printf("🔍 Block #%d Approval: %d/%d validators approved\n", block.Index, approvedVotes, len(pod.Validators))

	// ✅ Block is valid if 75% of validators approve
	if approvedVotes >= requiredVotes {
		fmt.Println("✅ Block approved by validators!")

		// ✅ Reward validators who approved the block
		for _, validator := range pod.Validators {
			if validator.ApproveBlock(block) {
				validator.RewardValidator(10) // Reward each approving validator
			}
		}

		return true
	}

	fmt.Println("❌ Block rejected due to insufficient votes.")
	return false
}
