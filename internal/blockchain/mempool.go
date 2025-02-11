package blockchain

import (
	"fmt"
	"sync"
)

// Mempool stores pending transactions before they are mined
type Mempool struct {
	Transactions []Transaction // List of transactions in mempool
	mu          sync.Mutex     // Mutex to prevent concurrent modification issues
}

// NewMempool initializes an empty transaction mempool
func NewMempool() *Mempool {
	return &Mempool{
		Transactions: []Transaction{},
	}
}

// AddTransaction adds a new transaction to the mempool
func (m *Mempool) AddTransaction(tx Transaction) {
	m.mu.Lock()         // Lock to prevent race conditions
	defer m.mu.Unlock() // Unlock after function execution

	m.Transactions = append(m.Transactions, tx)
	fmt.Printf("✅ Transaction added to mempool: %s\n", tx.TxID)
}

// GetTransactions returns all transactions in the mempool
func (m *Mempool) GetTransactions() []Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.Transactions // Return a copy of transactions
}

// Clear clears the mempool after mining a block
func (m *Mempool) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Transactions = []Transaction{} // Reset mempool
	fmt.Println("✅ Mempool cleared after block mining.")
}
