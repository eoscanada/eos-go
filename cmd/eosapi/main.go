package main

import (
	"fmt"
	"log"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
)

func main() {
	//api := eos.New(&url.URL{Scheme: "http", Host: "cbillett.eoscanada.com"}, bytes.Repeat([]byte{0}, 32))
	//api := eos.New("http://35.203.101.218:8888", bytes.Repeat([]byte{0}, 32))
	api := eos.New("http://localhost:8888")

	api.Debug = true
	eos.Debug = true

	keyBag := eos.NewKeyBag()
	for _, key := range []string{
		"5KEAGZjeSbWBoUHJLSZwsD5tWJp3JevXrczXGDi54zQYVy6C9HB",
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

		system.NewNewAccount(AC("eosio"), AC("aabbccddeeff"), ecc.MustNewPublicKey("EOS71UgDXVDXd56UUhCCLyn2U8QbjniESZLxRsBkkogNUHZwizY3b")),
		system.NewBuyRAM(
			eos.AccountName("eosio"),
			eos.AccountName("aabbccddeeff"),
			8192,
		),
		//system.NewDelegateBW(
		//	eos.AccountName("eosio"),
		//	eos.AccountName("bbbbbbbbbbbb"),
		//	eos.NewEOSAsset(10000),
		//	eos.NewEOSAsset(10000),
		//	true,
		//),
		//system.NewNewAccount(AC("eosio"), AC("aaaaaaaaaaac"), ecc.MustNewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")),

		//token.NewTransfer(eos.AccountName("eosio"), eos.AccountName("bbbbbbbbbbbb"), eos.NewEOSAsset(100000), ""),
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
