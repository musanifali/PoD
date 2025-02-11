package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"my_blockchain/internal/blockchain"
	"my_blockchain/api/routes"
)

// APIServer represents the HTTP API server
type APIServer struct {
	Port       string
	Blockchain *blockchain.Blockchain
}

// NewAPIServer initializes a new API server
func NewAPIServer(port string, blockchain *blockchain.Blockchain) *APIServer {
	return &APIServer{
		Port:       port,
		Blockchain: blockchain,
	}
}

// Start runs the API server
func (s *APIServer) Start() {
	router := mux.NewRouter()

	// Blockchain API Routes
	router.HandleFunc("/blocks", routes.GetBlocks(s.Blockchain)).Methods("GET")
	router.HandleFunc("/mine_block", routes.MineBlock(s.Blockchain)).Methods("POST")

	// Transaction Routes
	router.HandleFunc("/transactions", routes.GetTransactions(s.Blockchain)).Methods("GET")
	router.HandleFunc("/upload_file", routes.UploadFile(s.Blockchain)).Methods("POST")

	// Validator Routes
	router.HandleFunc("/validators", routes.GetValidators(s.Blockchain)).Methods("GET")
	router.HandleFunc("/register_validator", routes.RegisterValidator(s.Blockchain)).Methods("POST")

	// Wallet Routes
	router.HandleFunc("/wallet/create", routes.CreateWallet()).Methods("POST")
	router.HandleFunc("/wallet/sign", routes.SignData()).Methods("POST")

	// P2P Routes
	router.HandleFunc("/start_peer", routes.StartPeer(s.Blockchain)).Methods("POST")
	router.HandleFunc("/connect_peer", routes.ConnectPeer(s.Blockchain)).Methods("POST")
	router.HandleFunc("/sync_blockchain", routes.SyncBlockchain(s.Blockchain)).Methods("POST")

	// Start P2P server in the background
	go s.Blockchain.Network.StartServer()

	fmt.Println("🌐 API Server running on port", s.Port)
	log.Fatal(http.ListenAndServe(":"+s.Port, router))
}
