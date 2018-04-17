package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

func main() {
	//api := eos.New("http://testnet-dawn3.eosio.ca", "0000000000000000000000000000000000000000000000000000000000000000")
	api := eos.New(&url.URL{Scheme: "http", Host: "localhost:8889"}, bytes.Repeat([]byte{0}, 32))

	api.Debug = true

	keyBag := eos.NewKeyBag()
	if err := keyBag.Add("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"); err != nil {
		log.Fatalln("Couldn't load private key:", err)
	}
	//"EOS5GNc1NNsChC2URSevTBYhvQxGzjAgFRDi8V992ckjyc5tRzWn4",
	// "5K7Ffo8LXHhbsxV48w3sZzo8UnaKX3z5iD5mvac1AfDhHXKs3ao"
	if err := keyBag.Add("5K7Ffo8LXHhbsxV48w3sZzo8UnaKX3z5iD5mvac1AfDhHXKs3ao"); err != nil {
		log.Fatalln("Couldn't load private key:", err)
	}

	api.SetSigner(keyBag)

	walletAPI := eos.New(&url.URL{Scheme: "http", Host: "localhost:6667"}, bytes.Repeat([]byte{0}, 32))
	api.SetSigner(eos.NewWalletSigner(walletAPI, "default"))
	// Corresponding to the wallet, so we can sign on the live node.

	// resp, err := api.SetCode(AC("eosio"), "/home/abourget/build/eos/build/contracts/eosio.system/eosio.system.wasm", "/home/abourget/build/eos/build/contracts/eosio.system/eosio.system.abi")
	// if err != nil {
	// 	fmt.Println("ERROR calling SetCode:", err)
	// } else {
	// 	fmt.Println("RESP:", resp)
	// }

	resp, err := api.NewAccount(AC("eosio"), AC("abourget4"), ecc.MustNewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"))
	if err != nil {
		fmt.Println("ERROR calling NewAccount:", err)
	} else {
		fmt.Println("RESP:", resp)
	}

}
