package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"my_blockchain/internal/blockchain"
)

// ============================
// 🚀 Validator Routes
// ============================

// GetValidators returns a list of registered validators (public data only)
func GetValidators(bc *blockchain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// ✅ Filter validators to only return public data
		var publicValidators []map[string]interface{}
		for _, v := range bc.Consensus.Validators {
			publicValidators = append(publicValidators, map[string]interface{}{
				"ID":        v.ID,
				"PublicKey": v.PublicKey,
				"Balance":   v.Balance,
			})
		}

		fmt.Printf("✅ Returning %d registered validator(s).\n", len(publicValidators))
		json.NewEncoder(w).Encode(publicValidators)
	}
}

// RegisterValidator allows a new validator to join the network
func RegisterValidator(bc *blockchain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request struct {
			ID string `json:"ID"`
		}

		// ✅ Validate incoming request
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// ✅ Prevent duplicate validator registration
		for _, v := range bc.Consensus.Validators {
			if v.ID == request.ID {
				http.Error(w, "Validator already exists", http.StatusConflict)
				return
			}
		}

		// ✅ Create new validator with Public Key
		newValidator := blockchain.NewValidator(request.ID)
		bc.Consensus.RegisterValidator(newValidator)

		fmt.Printf("✅ Validator Registered: %s | Public Key: %s\n", newValidator.ID, newValidator.PublicKey)

		// ✅ Return structured response
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Validator registered successfully",
			"validator": map[string]interface{}{
				"ID":        newValidator.ID,
				"PublicKey": newValidator.PublicKey,
				"Balance":   newValidator.Balance,
			},
		})
	}
}
