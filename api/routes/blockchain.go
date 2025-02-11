package routes

import (
	"encoding/json"
	"net/http"
	"my_blockchain/internal/blockchain"
	"fmt"
	"math/rand"
)

// ============================
// 🚀 Blockchain Routes
// ============================

// GetBlocks returns all blocks in the blockchain
func GetBlocks(bc *blockchain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bc.Chain)
	}
}

// MineBlock mines a new block, ensuring validator approval using PoD
func MineBlock(bc *blockchain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// ✅ Debugging: Check Mempool Before Mining
		mempoolTransactions := bc.Mempool.GetTransactions()
		fmt.Printf("🔍 Mempool Before Mining: %v\n", mempoolTransactions)

		if len(mempoolTransactions) == 0 {
			http.Error(w, "❌ No transactions to mine!", http.StatusBadRequest)
			return
		}

		// ✅ Select a validator randomly to validate this block
		if len(bc.Consensus.Validators) == 0 {
			http.Error(w, "❌ No registered validators! Cannot mine block.", http.StatusInternalServerError)
			return
		}
		selectedValidator := bc.Consensus.Validators[rand.Intn(len(bc.Consensus.Validators))]
		fmt.Printf("🔍 Selected Validator for Mining: %s\n", selectedValidator.ID)

		// ✅ Mine transactions from the mempool into a block
		newBlock := bc.MineBlock(blockchain.NewWallet())

		if newBlock == nil {
			http.Error(w, "❌ Block mining failed! Not added to blockchain.", http.StatusInternalServerError)
			return
		}

		// ✅ Ensure validator approval before adding the block
		if !selectedValidator.ApproveBlock(*newBlock) {
			http.Error(w, "❌ Block rejected by validator!", http.StatusForbidden)
			return
		}

		fmt.Printf("✅ Block #%d successfully mined and added to the blockchain!\n", newBlock.Index)

		// ✅ Send response with ONLY the latest block (not full blockchain)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newBlock)
	}
}
