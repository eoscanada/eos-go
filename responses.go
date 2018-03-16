package eosapi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	"github.com/lunixbochs/struc"
)

type InfoResp struct {
	ServerVersion            string      `json:"server_version"`              // "2cc40a4e"
	HeadBlockNum             uint32      `json:"head_block_num"`              // 2465669,
	LastIrreversibleBlockNum uint32      `json:"last_irreversible_block_num"` // 2465655
	HeadBlockID              string      `json:"head_block_id"`               // "00259f856bfa142d1d60aff77e70f0c4f3eab30789e9539d2684f9f8758f1b88",
	HeadBlockTime            JSONTime    `json:"head_block_time"`             //  "2018-02-02T04:19:32"
	HeadBlockProducer        AccountName `json:"head_block_producer"`         // "inita"
	RecentSlots              string      `json:"recent_slots"`                //  "1111111111111111111111111111111111111111111111111111111111111111"
	ParticipationRate        string      `json:"participation_rate"`          // "1.00000000000000000" // this should be a `double`, or a decimal of some sort..

}

type BlockResp struct {
	Previous              string           `json:"previous"`                // : "0000007a9dde66f1666089891e316ac4cb0c47af427ae97f93f36a4f1159a194",
	Timestamp             JSONTime         `json:"timestamp"`               // : "2017-12-04T17:12:08",
	TransactionMerkleRoot string           `json:"transaction_merkle_root"` // : "0000000000000000000000000000000000000000000000000000000000000000",
	Producer              AccountName      `json:"producer"`                // : "initj",
	ProducerChanges       []ProducerChange `json:"producer_changes"`        // : [],
	ProducerSignature     string           `json:"producer_signature"`      // : "203dbf00b0968bfc47a8b749bbfdb91f8362b27c3e148a8a3c2e92f42ec55e9baa45d526412c8a2fc0dd35b484e4262e734bea49000c6f9c8dbac3d8861c1386c0",
	Cycles                []Cycle          `json:"cycles"`                  // : [],
	ID                    string           `json:"id"`                      // : "0000007b677719bdd76d729c3ac36bed5790d5548aadc26804489e5e179f4a5b",
	BlockNum              uint64           `json:"block_num"`               // : 123,
	RefBlockPrefix        uint64           `json:"ref_block_prefix"`        // : 2624744919

}

type GetTableRowsResp struct {
	More bool            `json:"more"`
	Rows json.RawMessage `json:"rows"` // defer loading, as it depends on `JSON` being true/false.
}

func (resp *GetTableRowsResp) JSONToStructs(v interface{}) error {
	return json.Unmarshal(resp.Rows, v)
}

func (resp *GetTableRowsResp) BinaryToStructs(v interface{}) error {
	var rows []string

	err := json.Unmarshal(resp.Rows, &rows)
	if err != nil {
		return err
	}

	for _, row := range rows {
		bin, err := hex.DecodeString(row)
		if err != nil {
			return err
		}

		ourstruct := &MyStruct{}
		if err := struc.Unpack(bytes.NewReader(bin), ourstruct); err != nil {
			return err
		}

		spew.Dump(ourstruct)
	}

	return nil
}

// type MyStruct struct {
// 	Key      string `struc:"[8]int8,little"`
// 	Balance  uint64 `struc:"uint64,little"`
// 	Currency string `struc:"[8]int8,little"`
// }

type MyStruct struct {
	Currency
	Balance uint64 `struc:"uint64,little"`
}

type Currency struct {
	Precision byte   `struc:"uint8"`
	Name      string `struc:"[7]uint8"`
}

type GetRequiredKeysResp struct {
	RequiredKeys []PublicKey `json:"required_keys"`
}

type PushTransactionResp struct {
	TransactionID string `json:"transaction_id"`
	Processed     bool   `json:"processed"` // WARN: is an `fc::variant` in server..
}
