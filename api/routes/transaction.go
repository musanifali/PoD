package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"my_blockchain/internal/blockchain"
)

// ============================
// 🚀 Transaction Routes
// ============================

// GetTransactions returns all transactions currently in the mempool (not yet mined)
func GetTransactions(bc *blockchain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		transactions := bc.Mempool.GetTransactions() // ✅ Fetch transactions from mempool
		if len(transactions) == 0 {
			http.Error(w, "No pending transactions", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(transactions)
	}
}

// UploadFile handles file uploads and creates a transaction
func UploadFile(bc *blockchain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Invalid file upload", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Ensure uploads directory exists
		uploadDir := "uploads/"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		// Save file locally
		filePath := uploadDir + header.Filename
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Failed to write file to disk", http.StatusInternalServerError)
			return
		}

		// Generate SHA-256 hash
		fileHash, err := blockchain.HashFile(filePath)
		if err != nil {
			http.Error(w, "Failed to generate file hash", http.StatusInternalServerError)
			return
		}

		// Create transaction
		wallet := blockchain.NewWallet()
		signature, err := wallet.SignData(fileHash)
		if err != nil {
			http.Error(w, "Failed to sign transaction", http.StatusInternalServerError)
			return
		}
		tx := blockchain.NewTransaction(fileHash, fmt.Sprintf("%x", wallet.PublicKey), header.Size, 0.0, signature)

		// ✅ Add transaction to the mempool instead of directly adding it to a block
		bc.AddTransaction(tx)

		fmt.Printf("✅ Transaction added to mempool: %+v\n", tx)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tx)
	}
}

