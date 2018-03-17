package main

import (
	"fmt"
	"log"

	"github.com/eosioca/eosapi"
)

func main() {
	api := eosapi.New("http://testnet-dawn3.eosio.ca")

	keybag := eosapi.NewKeyBag([]byte("hello-world-get-this-from-get-info-and-hex-decode-it"))
	if err := keybag.Add("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"); err != nil {
		log.Fatalln("Couldn't load private key:", err)
	}
	if err := keybag.Add("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss"); err != nil {
		log.Fatalln("Couldn't load private key:", err)
	}

	// Corresponding to the wallet, so we can sign on the live node.

	resp, err := api.SetCode(eosapi.AccountName("currency"), "file1.wasm", "file1.abi", keybag)
	if err != nil {
		fmt.Println("ERROR calling SetCode:", err)
	} else {
		fmt.Println("RESP:", resp)
	}
}
