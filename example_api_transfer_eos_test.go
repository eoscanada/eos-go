package eos_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/token"
)

func ExampleAPI_PushTransaction_transfer_EOS() {
	api := eos.New(getAPIURL())

	keyBag := &eos.KeyBag{}
	err := keyBag.ImportPrivateKey(context.Background(), readPrivateKey())
	if err != nil {
		panic(fmt.Errorf("import private key: %w", err))
	}
	api.SetSigner(keyBag)

	from := eos.AccountName("eosuser1")
	to := eos.AccountName("eosuser2")
	quantity, err := eos.NewEOSAssetFromString("1.0000 EOS")
	memo := ""

	if err != nil {
		panic(fmt.Errorf("invalid quantity: %w", err))
	}

	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(context.Background(), api); err != nil {
		panic(fmt.Errorf("filling tx opts: %w", err))
	}

	tx := eos.NewTransaction([]*eos.Action{token.NewTransfer(from, to, quantity, memo)}, txOpts)
	signedTx, packedTx, err := api.SignTransaction(context.Background(), tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		panic(fmt.Errorf("sign transaction: %w", err))
	}

	content, err := json.MarshalIndent(signedTx, "", "  ")
	if err != nil {
		panic(fmt.Errorf("json marshalling transaction: %w", err))
	}

	fmt.Println(string(content))
	fmt.Println()

	response, err := api.PushTransaction(context.Background(), packedTx)
	if err != nil {
		panic(fmt.Errorf("push transaction: %w", err))
	}

	fmt.Printf("Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
}

func readPrivateKey() string {
	// Right now, the key is read from an environment variable, it's an example after all.
	// In a real-world scenario, would you probably integrate with a real wallet or something similar
	envName := "EOS_GO_PRIVATE_KEY"
	privateKey := os.Getenv(envName)
	if privateKey == "" {
		panic(fmt.Errorf("private key environment variable %q must be set", envName))
	}

	return privateKey
}
