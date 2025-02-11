package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Block represents a single block in the blockchain
type Block struct {
	Index        int            // Block position in the chain
	Timestamp    string         // Block creation timestamp
	Transactions []Transaction  // Transactions stored in the block
	PreviousHash string         // Hash of the previous block
	Hash         string         // Unique block hash
	Signature    string         // Digital signature for authenticity
}

// NewBlock creates a new block containing validated transactions
func NewBlock(index int, transactions []Transaction, previousHash string, wallet *Wallet) Block {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	hash := calculateHash(index, timestamp, transactions, previousHash)

	// Sign the block hash using the wallet
	signature, _ := wallet.SignData(hash)

	return Block{
		Index:        index,
		Timestamp:    timestamp,
		Transactions: transactions,
		PreviousHash: previousHash,
		Hash:         hash,
		Signature:    signature,
	}
}

// calculateHash generates a SHA-256 hash for the block
func calculateHash(index int, timestamp string, transactions []Transaction, previousHash string) string {
	txID := "GENESIS" // Default TxID for empty transactions
	if len(transactions) > 0 {
		txID = transactions[0].TxID // Prevent index out of range
	}

	input := fmt.Sprintf("%d%s%s%s", index, timestamp, txID, previousHash)
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
