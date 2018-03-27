package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"

	"github.com/eosioca/eosapi"
	"github.com/eosioca/eosapi/ecc"
)

func main() {
	//api := eosapi.New("http://testnet-dawn3.eosio.ca", "0000000000000000000000000000000000000000000000000000000000000000")
	api := eosapi.New(&url.URL{Scheme: "http", Host: "localhost:18888"}, bytes.Repeat([]byte{0}, 32))

	keyBag := eosapi.NewKeyBag()
	if err := keyBag.Add("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"); err != nil {
		log.Fatalln("Couldn't load private key:", err)
	}
	if err := keyBag.Add("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss"); err != nil {
		log.Fatalln("Couldn't load private key:", err)
	}

	api.SetSigner(keyBag)

	walletAPI := eosapi.New(&url.URL{Scheme: "http", Host: "localhost:6667"}, bytes.Repeat([]byte{0}, 32))
	api.SetSigner(eosapi.NewWalletSigner(walletAPI))
	// Corresponding to the wallet, so we can sign on the live node.

	// resp, err := api.SetCode(eosapi.AccountName("currency"), "file1.wasm", "file1.abi", keybag)
	// if err != nil {
	// 	fmt.Println("ERROR calling SetCode:", err)
	// } else {
	// 	fmt.Println("RESP:", resp)
	// }

	resp, err := api.NewAccount(eosapi.AccountName("eosio"), eosapi.AccountName("abourget2"), ecc.MustNewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"))
	if err != nil {
		fmt.Println("ERROR calling NewAccount:", err)
	} else {
		fmt.Println("RESP:", resp)
	}

}
