package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Transaction represents a data upload or validation request
type Transaction struct {
	TxID       string   // Unique transaction ID
	FileHash   string   // SHA-256 hash of the file
	Uploader   string   // Public key of uploader
	Size       int64    // File size in bytes
	TrustScore float64  // AI-based file quality score
	Validators []string // List of validators who approved the transaction
	Signature  string   // Digital signature of uploader
}

// NewTransaction creates a new transaction for an uploaded file
func NewTransaction(fileHash string, uploader string, size int64, trustScore float64, signature string) Transaction {
	tx := Transaction{
		FileHash:   fileHash,
		Uploader:   uploader,
		Size:       size,
		TrustScore: 70,
		Signature:  signature,
	}

	// Generate a unique transaction ID (TxID)
	tx.TxID = tx.calculateTxID()
	return tx
}

// calculateTxID generates a SHA-256 hash as the transaction ID
func (tx *Transaction) calculateTxID() string {
	input := fmt.Sprintf("%s%s%d%f%s", tx.FileHash, tx.Uploader, tx.Size, tx.TrustScore, tx.Signature)
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
