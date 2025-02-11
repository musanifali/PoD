package main

import (
	"fmt"
	"os"
	"my_blockchain/api"
	"my_blockchain/internal/blockchain"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run api_main.go <port>")
		return
	}

	port := os.Args[1]

	// ✅ Initialize blockchain with the given port
	blockchain := blockchain.NewBlockchain(port)

	// ✅ Create and start the API server
	apiServer := api.NewAPIServer(port, blockchain)
	apiServer.Start()
}
