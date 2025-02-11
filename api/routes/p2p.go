package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"my_blockchain/internal/blockchain"
)

// StartPeer starts the node’s P2P server
func StartPeer(bc *blockchain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go bc.Network.StartServer()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "🌍 P2P Server started")
	}
}

// ConnectPeer allows a node to connect to another peer
func ConnectPeer(bc *blockchain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			PeerAddress string `json:"peer_address"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		bc.Network.ConnectToPeer(request.PeerAddress)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "🔗 Connected to peer: %s\n", request.PeerAddress)
	}
}

// SyncBlockchain allows nodes to request the latest blockchain data
func SyncBlockchain(bc *blockchain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bc.Network.BroadcastBlockchain()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "🔄 Blockchain synchronization requested")
	}
}
