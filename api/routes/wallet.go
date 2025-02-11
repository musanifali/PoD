package routes

import (
	"encoding/json"
	"net/http"
	"my_blockchain/internal/blockchain"
)

// CreateWallet generates a new wallet with cryptographic keys
func CreateWallet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wallet := blockchain.NewWallet()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(wallet)
	}
}

// SignData signs a given message using the wallet’s private key
func SignData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Message string `json:"message"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Create a wallet and sign the message
		wallet := blockchain.NewWallet()
		signature, _ := wallet.SignData(request.Message)

		// Send the signature as a response
		response := map[string]string{"signature": signature}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
