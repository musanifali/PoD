package blockchain

import (
	"fmt"
)

// Blockchain represents a list of blocks with PoD consensus, P2P networking, and a Mempool
type Blockchain struct {
	Chain     []Block       `json:"chain"`      // List of blocks in the blockchain
	Mempool   *Mempool      `json:"-"`          // ✅ Use separate mempool struct
	Consensus *PoDConsensus `json:"consensus"`  // Consensus mechanism
	Network   *P2PNetwork   `json:"-"`          // P2P network (excluded from JSON)
}

// NewBlockchain initializes the blockchain with PoD consensus, P2P networking, and an empty mempool
func NewBlockchain(port string) *Blockchain {
	genesisWallet := NewWallet()
	genesisBlock := NewBlock(0, []Transaction{}, "0", genesisWallet)
	pod := NewPoDConsensus()

	bc := &Blockchain{
		Chain:     []Block{genesisBlock},
		Mempool:   NewMempool(), // ✅ Use the new Mempool struct
		Consensus: pod,
	}

	// ✅ Ensure Network field is properly initialized only if not already connected
	if port != "" {
		bc.Network = NewP2PNetwork(bc, port)
	}

	return bc
}

// AddTransaction sends transactions to the mempool (NOT directly to the blockchain)
func (bc *Blockchain) AddTransaction(tx Transaction) {
	// ✅ Prevent duplicate transactions
	for _, existingTx := range bc.Mempool.GetTransactions() {
		if existingTx.TxID == tx.TxID {
			fmt.Println("⚠ Transaction already exists in mempool:", tx.TxID)
			return
		}
	}

	bc.Mempool.AddTransaction(tx) // ✅ Use Mempool struct instead of direct array
	fmt.Printf("✅ Transaction added to mempool: %s\n", tx.TxID)
}

// MineBlock moves transactions from mempool to a new block
func (bc *Blockchain) MineBlock(wallet *Wallet) *Block {
	transactions := bc.Mempool.GetTransactions() // ✅ Get transactions from mempool
	if len(transactions) == 0 {
		fmt.Println("⚠ No transactions in mempool to mine.")
		return nil
	}

	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := NewBlock(prevBlock.Index+1, transactions, prevBlock.Hash, wallet)

	// Validate and add block
	if !bc.Consensus.ValidateBlock(newBlock, bc) {
		fmt.Println("❌ Block validation failed! Not adding to blockchain.")
		return nil
	}

	bc.Chain = append(bc.Chain, newBlock)  // ✅ Add mined block to chain
	fmt.Printf("✅ Block #%d added with Proof-of-Data consensus!\n", newBlock.Index)

	bc.Mempool.Clear() // ✅ Clear mempool after block is added

	// ✅ Ensure `Network` is not nil before calling `BroadcastBlockchain`
	if bc.Network != nil {
		bc.Network.BroadcastBlockchain()
	} else {
		fmt.Println("⚠ Warning: Network is not initialized. Skipping broadcast.")
	}

	return &newBlock // ✅ Return newly mined block
}
