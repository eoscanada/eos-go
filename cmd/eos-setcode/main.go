package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
)

var flagAccount = flag.String("account", "", "Account to set code")
var flagWasm = flag.String("wasm", "", "WAST file loc")
var flagABI = flag.String("abi", "", "ABI file loc")
var flagPrivKey = flag.String("key", "", "Private key to load to sign transaction")

func main() {
	flag.Parse()

	api := eos.New(&url.URL{Scheme: "http", Host: "cbillett.eoscanada.com"}, bytes.Repeat([]byte{0}, 32))
	//api := eos.New(&url.URL{Scheme: "http", Host: "localhost:8889"}, bytes.Repeat([]byte{0}, 32))

	// api.Debug = true

	keyBag := eos.NewKeyBag()
	for _, key := range []string{
		*flagPrivKey,
		"5KE5hGNCAs1YvV74Ho14y1rV1DrnqZpTwLugS8QvYbKbrGAvVA1", // EOS71W8hvF43Eq6GQBRhuc5mvWKtknxzmb9NzNwPGpcEm2xAZaG8c
		"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3", //... 6CV
	} {
		if err := keyBag.Add(key); err != nil {
			log.Fatalln("Couldn't load private key:", err)
		}
	}

	api.SetSigner(keyBag)

	// WTTFFF, it DOESN'T WORK when we use `localhost:6667`, but it DOES work
	// when we use `localhost:5555` which is our GO wallet !!!
	// Then, it DOESN't work when we sign our transactions directly through
	// the KeyBag Signer
	// The only difference is that we GO THROUGH JSON SERIALIZATION and BACK ?!?
	walletAPI := eos.New(&url.URL{Scheme: "http", Host: "localhost:6667"}, bytes.Repeat([]byte{0}, 32))
	// walletAPI := eos.New(&url.URL{Scheme: "http", Host: "localhost:5555"}, bytes.Repeat([]byte{0}, 32))
	// walletAPI.Debug = true
	if err := walletAPI.WalletImportKey("default", *flagPrivKey); err != nil {
		fmt.Println("Error adding key to wallet:", err)
	}
	api.SetSigner(eos.NewWalletSigner(walletAPI, "default"))

	// Corresponding to the wallet, so we can sign on the live node.

	setCodeTx, err := system.NewSetCodeTx(eos.AccountName(*flagAccount), *flagWasm, *flagABI)
	if err != nil {
		log.Fatalln("Couldn't read setcode data:", err)
	}
	setCodeTx.Actions = []*eos.Action{setCodeTx.Actions[0]}

	resp, err := api.SignPushTransaction(setCodeTx, nil)
	if err != nil {
		fmt.Println("ERROR calling SetCode:", err)
	} else {
		fmt.Println("RESP:", resp)
	}

	api.SetSigner(keyBag)

	resp, err = api.SignPushTransaction(setCodeTx, nil)
	if err != nil {
		fmt.Println("ERROR calling SetCode:", err)
	} else {
		fmt.Println("RESP:", resp)
	}

}
