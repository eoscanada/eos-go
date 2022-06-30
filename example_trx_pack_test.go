package eos_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	eos "github.com/eoscanada/eos-go"
)

func ExamplePackedTransaction_Pack() {
	api := eos.New(getAPIURL())

	// Fills in the transaction header information like ref block, delays and expiration is set to 30s
	txOpts := &eos.TxOptions{}
	err := txOpts.FillFromChain(context.Background(), api)

	// The actual account here instead of `eosio` must be the the account that has the public key
	// associated with the private key below configure for 'active' level.
	from := "eosio"

	// This is of course now a burnt key, never ever use it, it's for demo purpose and protecting your key is your responsibility
	fromPrivateKey := "PVT_K1_2i6s2S8cxJw33zFuF3keAfUJjKSUJ53qVH7ac4veCuPVCpUSp"

	transferPermissionLevel, err := eos.NewPermissionLevel(from + "@active")
	NoError(err, "parse permission level")

	transferQuantity, err := eos.NewAssetFromString("10.0000 EOS")
	NoError(err, "parse asset")

	trx := eos.NewTransaction([]*eos.Action{
		{
			Account:       "eosio.token",
			Name:          "transfer",
			Authorization: []eos.PermissionLevel{transferPermissionLevel},
			ActionData: eos.NewActionData(&eos.Transfer{
				From:     eos.AccountName(from),
				To:       "eosio.token",
				Quantity: transferQuantity,
				Memo:     "Example action",
			}),
		},
	}, txOpts)

	keyBag := eos.NewKeyBag()
	err = keyBag.Add(fromPrivateKey)
	NoError(err, "add key to bag")

	keyBagKeys, err := keyBag.AvailableKeys(context.Background())
	NoError(err, "key bag available keys")
	Ensure(len(keyBagKeys) == 1, "expected a single available key")

	mainnetChainID := decodeChainID("aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906")
	signerKey := keyBagKeys[0]

	signedTrx, err := keyBag.Sign(context.Background(), eos.NewSignedTransaction(trx), mainnetChainID, signerKey)
	NoError(err, "sign transaction")

	packedTrx, err := signedTrx.Pack(eos.CompressionNone)
	NoError(err, "pack transaction")

	trxID, err := packedTrx.ID()
	NoError(err, "transaction id")

	fmt.Printf("Encode transaction: %s\n", encode(trx))
	fmt.Println()

	fmt.Printf("Encode signed transaction: %s\n", encode(signedTrx))
	fmt.Println()

	fmt.Printf("Encode packed and signed transaction ID %s: %s\n", trxID, encode(packedTrx))
	fmt.Println()
}

func encode(v interface{}) string {
	buffer := bytes.NewBuffer(nil)
	err := eos.NewEncoder(buffer).Encode(v)
	NoError(err, "encode %T", v)

	return hex.EncodeToString(buffer.Bytes())
}

func decodeChainID(in string) eos.Checksum256 {
	data, err := hex.DecodeString(in)
	NoError(err, "chain ID %q is not a valid checksum256 value", in)
	Ensure(len(data) == 32, "invalid checksum256, must have 32 bytes got %d", len(data))

	return data
}
