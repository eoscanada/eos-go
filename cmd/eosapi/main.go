package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
)

var flagPrivKey = flag.String("key", "", "The key to sign with")

func main() {
	flag.Parse()

	// api := eos.New(&url.URL{Scheme: "http", Host: "cbillett.eoscanada.com"}, bytes.Repeat([]byte{0}, 32))
	//api := eos.New(&url.URL{Scheme: "http", Host: "localhost:8889"}, bytes.Repeat([]byte{0}, 32))
	api := eos.New(&url.URL{Scheme: "http", Host: "blastx.eoscanada.com"}, bytes.Repeat([]byte{0}, 32))

	// api.Debug = true

	keyBag := eos.NewKeyBag()
	for _, key := range []string{
		*flagPrivKey,
		"5KaM7o9cAyAEH1kRvMQdEzWSv1LNwcpDGFhzajn8Fc8MB8rEtpY", // EOS74WXa58iuNnJaXVkh9AKKiHAxihyLtTAQtasTteLJHkj35JUNx
		"5J6rkLgUEgMThg6GH2iNva5RvRqmGabbhMFwvL6TvdXCp2cU7LR", // EOS6H7qNHBhD31wZ5k7CKavMJAfj2TxFgZcYxLhfSaTwknXuGBptE
		"5JgFcyKbovph99NNeUptnRRbyqGLmx9zB5sMdgQjPfzpWXotxSs", // EOS8Ci1VAkm6WxpGrVatfbub2sfqPtpc11jvuwcR8bdZwNam2JVdx
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

	newAcct := AC("abou2234")
	resp, err := api.SignPushActions(
		system.NewNewAccount(AC("eosio"), newAcct, ecc.MustNewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")),
	)
	if err != nil {
		fmt.Println("ERROR calling NewAccount:", err)
	} else {
		fmt.Println("RESP:", resp)
	}

	// walletAPI := eos.New(&url.URL{Scheme: "http", Host: "localhost:5555"}, bytes.Repeat([]byte{0}, 32))
	// walletAPI := eos.New(&url.URL{Scheme: "http", Host: "localhost:6667"}, bytes.Repeat([]byte{0}, 32))
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
