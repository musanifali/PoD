package blockchain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

// Peer-to-Peer Network
type P2PNetwork struct {
	Peers      []string
	Blockchain *Blockchain
	Port       string
}

// NewP2PNetwork initializes the P2P network
func NewP2PNetwork(blockchain *Blockchain, port string) *P2PNetwork {
	return &P2PNetwork{
		Peers:      []string{},
		Blockchain: blockchain,
		Port:       port,
	}
}


// StartServer starts a P2P server on the correct local IP
func (p2p *P2PNetwork) StartServer() {
	// Force binding to local IP instead of APIPA (169.254.x.x)
	localIP := getLocalIPv4() // Get correct local IP (192.168.x.x or 127.0.0.1)
	address := localIP + ":" + p2p.Port

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("❌ Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("🌍 P2P Server started on", address)
	fmt.Printf("🔗 Your Node Address: %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("❌ Connection error:", err)
			continue
		}
		go p2p.HandleConnection(conn)
	}
}


// getLocalIPv4 returns the local IPv4 address of the machine
// getLocalIPv4 returns the correct local IPv4 address (excludes APIPA)
func getLocalIPv4() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			// ✅ Ignore APIPA addresses (169.254.x.x)
			if strings.HasPrefix(ipNet.IP.String(), "169.254.") {
				continue
			}
			return ipNet.IP.String()
		}
	}

	return "127.0.0.1" // Default if no valid local IP found
}


// ConnectToPeer connects to a remote peer
func (p2p *P2PNetwork) ConnectToPeer(address string) {
	conn, err := net.Dial("tcp4", address) // Force IPv4
	if err != nil {
		fmt.Println("❌ Failed to connect to peer:", err)
		return
	}
	defer conn.Close()

	p2p.Peers = append(p2p.Peers, address)
	fmt.Println("🔗 Connected to peer:", address)

	// Send our blockchain to the new peer
	p2p.SendBlockchain(conn)
}

// HandleConnection processes incoming peer connections
func (p2p *P2PNetwork) HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("❌ Error reading data:", err)
		return
	}

	// Determine if we received a blockchain update
	if strings.HasPrefix(data, "BLOCKCHAIN_UPDATE") {
		var receivedChain Blockchain
		err = json.Unmarshal([]byte(strings.TrimPrefix(data, "BLOCKCHAIN_UPDATE ")), &receivedChain)
		if err != nil {
			fmt.Println("❌ Error decoding blockchain:", err)
			return
		}

		// Synchronize if the received chain is longer
		if len(receivedChain.Chain) > len(p2p.Blockchain.Chain) {
			fmt.Println("🔄 Synchronizing blockchain from peer...")
			p2p.Blockchain.Chain = receivedChain.Chain
		}
	}
}

// SendBlockchain sends our blockchain to a peer
func (p2p *P2PNetwork) SendBlockchain(conn net.Conn) {
	blockchainData, _ := json.Marshal(p2p.Blockchain)
	conn.Write([]byte("BLOCKCHAIN_UPDATE " + string(blockchainData) + "\n"))
}
// BroadcastBlockchain sends the blockchain to all connected peers
func (p2p *P2PNetwork) BroadcastBlockchain() {
	blockchainData, err := json.Marshal(p2p.Blockchain)
	if err != nil {
		fmt.Println("❌ Error encoding blockchain:", err)
		return
	}

	for _, peer := range p2p.Peers {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			fmt.Println("❌ Failed to connect to peer:", err)
			continue
		}
		defer conn.Close()

		conn.Write([]byte("BLOCKCHAIN_UPDATE " + string(blockchainData) + "\n"))
		fmt.Println("🔄 Sent blockchain update to", peer)
	}
}
