package eos

import "github.com/eoscanada/eos-go/ecc"

type HandshakeMessage struct {
	// net_plugin/protocol.hpp handshake_message
	NetworkVersion           int16         `json:"network_version"`
	ChainID                  HexBytes      `json:"chain_id"`
	NodeID                   HexBytes      `json:"node_id"` // sha256
	Key                      ecc.PublicKey `json:"key"`     // can be empty, producer key, or peer key
	Time                     Tstamp        `json:"time"`    // time?!
	Token                    HexBytes      `json:"token"`   // digest of time to prove we own the private `key`
	Signature                ecc.Signature `json:"sig"`     // can be empty if no key, signature of the digest above
	P2PAddress               string        `json:"p2p_address"`
	LastIrreversibleBlockNum uint32        `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  HexBytes      `json:"last_irreversible_block_id"`
	HeadNum                  uint32        `json:"head_num"`
	HeadID                   HexBytes      `json:"head_id"`
	OS                       string        `json:"os"`
	Agent                    string        `json:"agent"`
	Generation               int16         `json:"generation"`
}

type GoAwayReason uint8

const (
	GoAwayNoReason = uint8(iota)
	GoAwaySelfConnect
	GoAwayDuplicate
	GoAwayWrongChain
	GoAwayWrongVersion
	GoAwayForked
	GoAwayUnlinkable
	GoAwayBadTransaction
	GoAwayValidation
	GoAwayAuthentication
	GoAwayFatalOther
	GoAwayBenignOther
	GoAwayCrazy
)

type GoAwayMessage struct {
	GoAwayReason
}

type TimeMessage struct {
	Origin      Tstamp `json:"org"`
	Receive     Tstamp `json:"rec"`
	Transmit    Tstamp `json:"xmt"`
	Destination Tstamp `json:"dst"`
}

type SignedBlockSummaryMessage struct {
}
