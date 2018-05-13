package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
)

func main() {
	//api := eos.New(&url.URL{Scheme: "http", Host: "cbillett.eoscanada.com"}, bytes.Repeat([]byte{0}, 32))
	api := eos.New("http://localhost:8888", bytes.Repeat([]byte{0}, 32))
	//api := eos.New(&url.URL{Scheme: "http", Host: "localhost:8889"}, bytes.Repeat([]byte{0}, 32))

	api.Debug = true
	eos.Debug = true

	keyBag := eos.NewKeyBag()
	for _, key := range []string{
		"5JuWeC5KwZRVUQZ4eneYCYQ6Pa132QgvDQzEVJBA7XTgNTBWWRw",
		"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3",
		"5K6CAnUcDpJzBJBCve3QfQjEsxHrC8jqYgnYE1tDv4rKbDLG58N", // latest bios boot
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

	//resp, err := api.GetTransactions(eos.AccountName("cbillett"))
	//if err != nil {
	//	fmt.Println("Grr", err)
	//}
	//fmt.Println(resp)

	actionResp, err := api.SignPushActions(

		system.NewNewAccount(AC("eosio"), AC("aaaaaaaaaaab"), ecc.MustNewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")),

		//		token.NewTransfer(eos.AccountName("cbillett"), eos.AccountName("bozo"), eos.NewEOSAsset(100), ""),
	)
	if err != nil {
		fmt.Println("ERROR calling :", err)
	} else {
		fmt.Println("RESP:", actionResp)
	}

	// resp, err := api.GetCurrencyBalance(AC("eosio"), eos.EOSSymbol.Symbol, AC("eosio"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(resp)

}
