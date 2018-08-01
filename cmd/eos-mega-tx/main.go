package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/cmd/eos-mega-tx/random"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/eoscanada/eos-go/token"
	"github.com/satori/go.uuid"
)

var signingKey = flag.String("signing-key", "", "Key to sign transactions we're about to blast")
var apiAddr = flag.String("api-addr", "http://localhost:8888", "RPC endpoint of the nodeos instance")

var eosioAccount = AC("eosio")

type transferChannel chan *eos.Action

func main() {
	flag.Parse()

	privKey, err := ecc.NewPrivateKey(*signingKey)
	if err != nil {
		log.Fatalln("failed loading private key:", err)
	}

	api := eos.New(*apiAddr) // TODO: use chain ID from somewhere..

	keyBag := eos.NewKeyBag()
	if err := keyBag.Add(*signingKey); err != nil {
		log.Fatalln("Couldn't load private key:", err)
	}
	api.SetSigner(keyBag)

	var accountNames []eos.AccountName
	for i := 0; i < 10; i++ {
		acctName := random.String(12)

		fmt.Println("Creating account with name: ", acctName)

		newAcctAction := system.NewNewAccount(eosioAccount, AC(acctName), privKey.PublicKey())
		xferAction := token.NewTransfer(eosioAccount, AC(acctName), eos.NewEOSAsset(10000), "")
		uuid, err := uuid.NewV4()
		if err != nil {
			log.Fatal(uuid)
		}
		nonceAction := system.NewNonce(uuid.String())

		fmt.Printf("Will transfer from account : [%s] to account : [%s]\n", eosioAccount, acctName)
		_, err = api.SignPushActions(newAcctAction, xferAction, nonceAction)
		if err != nil {
			log.Fatalln("ERROR pushing tx:", err)
		}

		accountNames = append(accountNames, AC(acctName))
	}

	l := len(accountNames)
	fromIndex, toIndex := 0, 1

	for {
		from := accountNames[fromIndex]
		to := accountNames[toIndex]

		fmt.Printf("Transfer 0.1 EOS from [%d %s] to [%d %s]\n", fromIndex, from, toIndex, to)

		uuid, err := uuid.NewV4()
		if err != nil {
			log.Fatal(err)
		}

		_, err = api.SignPushActions(
			token.NewTransfer(from, to, eos.NewEOSAsset(1000), ""),

			system.NewNonce(uuid.String()),
		)
		if err != nil {
			log.Fatalln("ERROR sending transfer:", err)
		}

		fromIndex++
		toIndex++

		if fromIndex >= l {
			fromIndex = 0
		}
		if toIndex >= l {
			toIndex = 0
		}
	}
}
