package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/cmd/eos-mega-tx/random"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/eoscanada/eos-go/token"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var eosioAccount = AC("eosio")

type transferChannel chan *eos.Action

func main() {

	//api := eos.New("http://testnet-dawn3.eosio.ca", "0000000000000000000000000000000000000000000000000000000000000000")
	api := eos.New(&url.URL{Scheme: "http", Host: "localhost:18888"}, bytes.Repeat([]byte{0}, 32))

	keyBag := eos.NewKeyBag()
	if err := keyBag.Add("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"); err != nil {
		log.Fatalln("Couldn't load private key:", err)
	}
	if err := keyBag.Add("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss"); err != nil {
		log.Fatalln("Couldn't load private key:", err)
	}

	//Connect to proxy ...
	resp, err := api.NetConnect("192.168.1.147:29876")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connect to proxy reponse: ", resp)
	}

	api.SetSigner(keyBag)

	walletAPI := eos.New(&url.URL{Scheme: "http", Host: "localhost:16666"}, bytes.Repeat([]byte{0}, 32))
	api.SetSigner(eos.NewWalletSigner(walletAPI, "default"))

	accountActions, err := generateAccountActions(10)
	if err != nil {
		panic(err)
	}

	var accountNames []eos.AccountName
	for _, aa := range accountActions {

		resp, err := api.SignPushActions(aa)

		if err != nil {
			fmt.Println("ERROR calling NewAccount:", err)
			panic(err)
		} else {
			fmt.Println("Account Created:", resp)
		}

		a := aa.Obj().(system.NewAccount)
		accountNames = append(accountNames, a.Name)

		fmt.Printf("Will transfer form account : [%s] to account : [%s]\n", eosioAccount, a)
		resp, err = api.SignPushActions(token.NewTransfer(eosioAccount, a.Name, eos.NewEOSAsset(10000), ""))
		if err != nil {
			fmt.Println("ERROR transfering: ", err)

		} else {
			fmt.Println("Transfer: ", resp)
		}

	}

	c := make(transferChannel, 10)

	for i := 0; i < 1; i++ {
		go sendTransfer(c, api)
	}

	done := make(chan bool)
	go generateTransfer(c, accountNames)

	<-done
}

func sendTransfer(c transferChannel, api *eos.API) error {

	fmt.Println("Setting new transfer sender ...")

	for a := range c {

		t := a.Obj().(token.Transfer)
		uuid, _ := uuid.NewV4()

		resp, err := api.SignPushActions(a, system.NewNonce(uuid.String()))
		if err != nil {
			fmt.Printf("ERROR transfering from [%s] to [%s] error [%s]\n", t.From, t.To, err)

		} else {
			fmt.Println("Transfer: ", resp)
		}
	}

	return nil
}

func generateTransfer(c transferChannel, accountNames []eos.AccountName) error {

	fmt.Println("Starting transfer generation")

	l := len(accountNames)

	if l < 2 {
		return errors.New("Expecting at least 2 account names")
	}

	reset := func() (int, int) { return 0, 1 }
	fromIndex, toIndex := reset()

	go func() {
		for {
			from := accountNames[fromIndex]
			to := accountNames[toIndex]

			fmt.Printf("Adding transfert [%d %s-> %d %s] to channel\n", fromIndex, from, toIndex, to)
			c <- token.NewTransfer(from, to, eos.NewEOSAsset(1000), "")

			fromIndex += 1
			toIndex += 1
			if fromIndex >= l {
				fromIndex = 0
			}
			if toIndex >= l {
				toIndex = 0
			}
		}
	}()

	return nil
}

func generateAccountActions(count int) (actions []*eos.Action, err error) {

	for i := 0; i < count; i++ {

		name := random.String(12)
		fmt.Println("Creating account with name: ", name)
		if err != nil {
			panic(err)
		}

		a := system.NewNewAccount(eosioAccount, AC(name), ecc.MustNewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"))
		actions = append(actions, a)
	}

	return
}
