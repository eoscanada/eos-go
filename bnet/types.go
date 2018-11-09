package bnet

import (
	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type BNetMessageType byte

const (
	HelloType                    BNetMessageType = iota // 0
	TransactionNoticeType                               // 1
	BlockNoticeType                                     // 2
	SignedBlockPointerType                              // 3
	PackedTransactionPointerType                        // 4
	PingType                                            // 5
	PongType                                            // 6
)

//
// Message types
//

type Hello struct {
	PeerID                   ecc.PublicKey     `json:"peer_id"`
	NetworkVersion           string            `json:"network_version"`
	User                     string            `json:"user"`
	Password                 string            `json:"password"`
	Agent                    string            `json:"agent"`
	ProtocolVersion          string            `json:"protocol_version"`
	ChainID                  eos.Checksum256   `json:"chain_id"`
	RequestTransactions      eos.Bool          `json:"request_transactions"`
	LastIrreversibleBlockNum uint32            `json:"last_irr_block_num"`
	PendingBlockIDs          []eos.Checksum256 `json:"pending_block_ids"`
}

/**
 * From bnet_plugin.cpp's `trx_notice`:
 *
 * This message is sent upon successful speculative application of a transaction
 * and informs a peer not to send this message.
 */
type TransactionNotice struct {
	SignedTransactionIDs []eos.Checksum256 ///< hash of trx + sigs
}

/**
 * From bnet_plugin.cpp's `block_notice`:
 *
 * This message is sent upon successfully adding a transaction to the fork database
 * and informs the remote peer that there is no need to send this block.
 */
type BlockNotice struct {
	BlockIDs []eos.Checksum256 `json:"block_ids"`
}

type Ping struct {
	Sent                  eos.Tstamp      `json:"sent"`
	Code                  eos.Checksum256 `json:"code"`
	LastIrreversibleBlock uint32          `json:"lib"` /// last irreversible block
}

type Pong struct {
	Sent eos.Tstamp      `json:"sent"`
	Code eos.Checksum256 `json:"code"`
}

// Also use `eos.SignedBlock`
// Also use `eos.SignedTransaction`
var SignedBlock = eos.SignedTransaction{}
