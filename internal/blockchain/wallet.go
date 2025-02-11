package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
)

// Wallet represents a user's private and public keys
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

// NewWallet generates a new wallet with an ECDSA key pair
func NewWallet() *Wallet {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("❌ Error generating key pair:", err)
		return nil
	}
	return &Wallet{PrivateKey: privateKey, PublicKey: &privateKey.PublicKey}
}

// SignData signs a given hash using the private key
func (w *Wallet) SignData(data string) (string, error) {
	hash := sha256.Sum256([]byte(data))
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, hash[:])
	if err != nil {
		return "", err
	}

	// Convert r and s to base64 for easy storage
	rText := base64.StdEncoding.EncodeToString(r.Bytes())
	sText := base64.StdEncoding.EncodeToString(s.Bytes())

	return fmt.Sprintf("%s.%s", rText, sText), nil
}

// VerifySignature verifies a signed transaction or block
func VerifySignature(publicKey *ecdsa.PublicKey, data string, signature string) bool {
	hash := sha256.Sum256([]byte(data))

	// Split the signature into r and s components
	parts := strings.Split(signature, ".")
	if len(parts) != 2 {
		return false
	}

	rBytes, _ := base64.StdEncoding.DecodeString(parts[0])
	sBytes, _ := base64.StdEncoding.DecodeString(parts[1])

	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)

	return ecdsa.Verify(publicKey, hash[:], r, s)
}
