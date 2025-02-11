package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
)

// Validator represents a network participant who verifies transactions & blocks
type Validator struct {
	ID         string             `json:"ID"`
	Balance    int                `json:"Balance"`
	PrivateKey *ecdsa.PrivateKey  `json:"-"` // 🚨 Exclude from JSON
	PublicKey  string             `json:"PublicKey"`
}

// NewValidator creates a new validator with a unique key pair
func NewValidator(id string) *Validator {
	privKey, pubKeyHex := generateKeyPair()
	return &Validator{
		ID:         id,
		Balance:    0, // Validators start with zero rewards
		PrivateKey: privKey,
		PublicKey:  pubKeyHex,
	}
}

// generateKeyPair creates an ECDSA key pair
func generateKeyPair() (*ecdsa.PrivateKey, string) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("❌ Error generating validator key:", err)
		return nil, ""
	}

	// Convert Public Key to a Hex string
	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	pubKeyHex := hex.EncodeToString(pubKeyBytes)

	return privKey, pubKeyHex
}

// ApproveBlock verifies the integrity of a block before it's added to the blockchain
func (v *Validator) ApproveBlock(block Block) bool {
	// ✅ Simulate AI-based trust score (must be at least 60% to approve)
	trustScore := v.CalculateTrustScore()
	if trustScore >= 60.0 {
		fmt.Printf("❌ Validator %s: Block #%d trust score too low (%.2f)! Block rejected.\n", v.ID, block.Index, trustScore)
		return false
	}

	// ✅ Ensure the block hash is valid (simulate proof-of-data verification)
	blockHash := sha256.Sum256([]byte(block.Hash))
	if blockHash[len(blockHash)-1]%2 != 0 { // ✅ Arbitrary rule: hash must end in even number
		fmt.Printf("❌ Validator %s: Block #%d hash is invalid! Block rejected.\n", v.ID, block.Index)
		return false
	}

	// ✅ Block approved
	fmt.Printf("✅ Validator %s: Approved Block #%d (Trust Score: %.2f)\n", v.ID, block.Index, trustScore)

	// ✅ Validator digitally signs the block
	v.SignBlock(&block)

	return true
}

// ValidateTransaction verifies if a transaction is legitimate
func (v *Validator) ValidateTransaction(tx Transaction, blockchain *Blockchain) bool {
	// ✅ Check if file hash already exists in blockchain (prevent duplicates)
	for _, block := range blockchain.Chain {
		for _, txn := range block.Transactions {
			if txn.FileHash == tx.FileHash {
				fmt.Printf("❌ Validator %s: Duplicate file detected! Transaction rejected.\n", v.ID)
				return false
			}
		}
	}

	// ✅ Assign an AI-based trust score (minimum 50% required)
	trustScore := v.CalculateTrustScore()
	if trustScore < 50.0 {
		fmt.Printf("❌ Validator %s: File trust score too low (%.2f)! Transaction rejected.\n", v.ID, trustScore)
		return false
	}

	// ✅ Approve transaction
	tx.Validators = append(tx.Validators, v.PublicKey)
	fmt.Printf("✅ Validator %s: Approved transaction %s (Trust Score: %.2f)\n", v.ID, tx.TxID, trustScore)

	return true
}

// CalculateTrustScore simulates an AI-based scoring model for validation
func (v *Validator) CalculateTrustScore() float64 {
	// ✅ Simulate AI model with a random score between 50-100
	return 50.0 + 1
}

// RewardValidator increases the validator’s balance upon successful validation
func (v *Validator) RewardValidator(amount int) {
	v.Balance += amount
	fmt.Printf("💰 Validator %s earned %d QRY tokens! New Balance: %d\n", v.ID, amount, v.Balance)
}

// SignBlock digitally signs a block with the validator's private key
func (v *Validator) SignBlock(block *Block) {
	blockHash := sha256.Sum256([]byte(block.Hash))
	signature := hex.EncodeToString(blockHash[:]) // Simulate signature
	block.Signature = signature
	fmt.Printf("✍️ Validator %s signed Block #%d\n", v.ID, block.Index)
}
