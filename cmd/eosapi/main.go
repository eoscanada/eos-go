package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
)

func main() {
	// api := eos.New(&url.URL{Scheme: "http", Host: "cbillett.eoscanada.com"}, bytes.Repeat([]byte{0}, 32))
	api := eos.New(&url.URL{Scheme: "http", Host: "localhost:8889"}, bytes.Repeat([]byte{0}, 32))

	// api.Debug = true

	keyBag := eos.NewKeyBag()
	for _, key := range []string{
		"5KE5hGNCAs1YvV74Ho14y1rV1DrnqZpTwLugS8QvYbKbrGAvVA1", // EOS71W8hvF43Eq6GQBRhuc5mvWKtknxzmb9NzNwPGpcEm2xAZaG8c
		"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3", //... 6CV
		"5Jrwky4GxChTSqG29Mj9B1HGqJXx8T8WxkPJULmDaBDsguhiF8m",
	} {
		if err := keyBag.Add(key); err != nil {
			log.Fatalln("Couldn't load private key:", err)
		}
	}

	api.SetSigner(keyBag)

	// Corresponding to the wallet, so we can sign on the live node.

	// resp, err := api.SetCode(AC("eosio"), "/home/abourget/build/eos/build/contracts/eosio.system/eosio.system.wasm", "/home/abourget/build/eos/build/contracts/eosio.system/eosio.system.abi")
	// if err != nil {
	// 	fmt.Println("ERROR calling SetCode:", err)
	// } else {
	// 	fmt.Println("RESP:", resp)
	// }

	newAcct := AC("abou23u")

	resp, err := api.SignPushActions(
		system.NewNewAccount(AC("eosio"), newAcct, ecc.MustNewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")),
	)
	if err != nil {
		fmt.Println("ERROR calling NewAccount:", err)
	} else {
		fmt.Println("RESP:", resp)
	}

	// walletAPI := eos.New(&url.URL{Scheme: "http", Host: "localhost:5555"}, bytes.Repeat([]byte{0}, 32))
	// // walletAPI.Debug = true
	// api.SetSigner(eos.NewWalletSigner(walletAPI, "default"))

	// resp, err = api.SignPushActions(
	// 	system.NewNewAccount(AC("eosio"), newAcct, ecc.MustNewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")),
	// )
	// if err != nil {
	// 	fmt.Println("ERROR calling NewAccount:", err)
	// } else {
	// 	fmt.Println("RESP:", resp)
	// }
}
