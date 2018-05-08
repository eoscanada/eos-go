package main

import (
	"bytes"
	"log"
	"net/url"

	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/token"
)

func main() {
	//api := eos.New(&url.URL{Scheme: "http", Host: "cbillett.eoscanada.com"}, bytes.Repeat([]byte{0}, 32))
	api := eos.New(&url.URL{Scheme: "http", Host: "Charless-MacBook-Pro-2.local:8888"}, bytes.Repeat([]byte{0}, 32))
	//api := eos.New(&url.URL{Scheme: "http", Host: "localhost:8889"}, bytes.Repeat([]byte{0}, 32))

	api.Debug = true

	keyBag := eos.NewKeyBag()
	for _, key := range []string{
		"5J5EE2cBDM4d3vWpKGcJsgiagsLVZkgWjJpxacz9mXodemXex6K",
		"5Jd9CCuMGENFJTk1RGiCWCtLhCzkHcDLBnc8vnhGMArFu5dBfYF",
		"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3",
		"5J77j8KYX33cgVPMQZ82zD967VNA9SPcXWnjRkb27z9M2suaZNn",
		"5JJbFqMRLncsRXbVYSUwdMyQke1ULLH65nBLBsDPnxARDdsYnhK",
	} {
		if err := keyBag.Add(key); err != nil {
			log.Fatalln("Couldn't load private key:", err)
		}
	}

	api.SetSigner(keyBag)

	// Corresponding to the wallet, so we can sign on the live node.

	//setCodeTx, err := system.NewSetCodeTx(
	//	AC("eosio"),
	//	"/Users/cbillett/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.system.wasm",
	//	"/Users/cbillett/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.system.abi",
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//resp, err := api.SignPushTransaction(setCodeTx, &eos.TxOptions{})
	//if err != nil {
	//	fmt.Println("ERROR calling NewAccount:", err)
	//} else {
	//	fmt.Println("RESP:", resp)
	//}

	//setCodeTx, err := system.NewSetCodeTx(
	//	AC("eosio"),
	//	"/Users/cbillett/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.bios.wasm",
	//	"/Users/cbillett/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.bios.abi",
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	////
	//resp, err := api.SignPushTransaction(setCodeTx, &eos.TxOptions{})
	//if err != nil {
	//	fmt.Println("ERROR calling NewAccount:", err)
	//} else {
	//	fmt.Println("RESP:", resp)
	//}

	//setCodeTx, err := system.NewSetCodeTx(
	//	AC("eosio.msig"),
	//	"/Users/cbillett/devel/dix975/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.msig.wasm",
	//	"/Users/cbillett/devel/dix975/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.msig.abi",
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//resp, err := api.SignPushTransaction(setCodeTx, &eos.TxOptions{})
	//if err != nil {
	//	fmt.Println("ERROR calling NewAccount:", err)
	//} else {
	//	fmt.Println("RESP:", resp)
	//}

	//setCodeTx, err := system.NewSetCodeTx(
	//	AC("eosio.token"),
	//	"/Users/cbillett/devel/dix975/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.token.wasm",
	//	"/Users/cbillett/devel/dix975/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.token.abi",
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//resp, err := api.SignPushTransaction(setCodeTx, &eos.TxOptions{})
	//if err != nil {
	//	fmt.Println("ERROR calling NewAccount:", err)
	//} else {
	//	fmt.Println("RESP:", resp)
	//}

	//resp, err := api.GetTransaction("a343fcb285088bc955f721b9f22efad3e2fd131bad93421364f1b043a3aad00f")
	//if err != nil {
	//	fmt.Println("Grr", err)
	//}
	//fmt.Println(resp)

	actionResp, err := api.SignPushActions(

		//system.NewNewAccount(AC("eosio"), AC("cbillett"), ecc.MustNewPublicKey("EOS66MfGpiepzs46DudrpSQw6GEn2QywFYVMWc18hBFVVVehdbKdi")),

		token.NewTransfer(eos.AccountName("eosio"), eos.AccountName("cbillett"), eos.NewEOSAsset(100000), ""),
	)
	if err != nil {
		fmt.Println("ERROR calling :", err)
	} else {
		fmt.Println("RESP:", actionResp)
	}

	resp, err := api.GetCurrencyBalance(AC("eosio"), eos.EOSSymbol.Symbol, AC("eosio"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)

}
