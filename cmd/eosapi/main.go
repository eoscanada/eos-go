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
	api := eos.New(&url.URL{Scheme: "http", Host: "Charless-MacBook-Pro-2.local:18888"}, bytes.Repeat([]byte{0}, 32))
	//api := eos.New(&url.URL{Scheme: "http", Host: "localhost:8889"}, bytes.Repeat([]byte{0}, 32))

	//api.Debug = true

	keyBag := eos.NewKeyBag()
	for _, key := range []string{
		"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3",
		"5J77j8KYX33cgVPMQZ82zD967VNA9SPcXWnjRkb27z9M2suaZNn",
		"5JJbFqMRLncsRXbVYSUwdMyQke1ULLH65nBLBsDPnxARDdsYnhK",
	} {
		if err := keyBag.Add(key); err != nil {
			log.Fatalln("Couldn't load private key:", err)
		}
	}

	api.SetSigner(keyBag)

	var err error

	// Corresponding to the wallet, so we can sign on the live node.

	//setCodeTx, err := system.NewSetCodeTx(
	//	AC("eosio"),
	//	"/Users/cbillett/devel/dix975/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.system.wasm",
	//	"/Users/cbillett/devel/dix975/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.system.abi",
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
	//	"/Users/cbillett/devel/dix975/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.bios.wasm",
	//	"/Users/cbillett/devel/dix975/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.bios.abi",
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
	//	AC("eosio"),
	//	"/Users/cbillett/devel/dix975/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.system.wasm",
	//	"/Users/cbillett/devel/dix975/go/src/github.com/diagramventures/bc/bios-docker/contracts/eosio.system.abi",
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

	actionResp, err := api.SignPushActions(

		//system.NewNewAccount(AC("eosio"), AC("eosio.msig"), ecc.MustNewPublicKey("EOS8ju96GnaKaYAs7b5EvvtwWqTVPepSCciDHvDCjiEhGb5joYtjk")),
		//system.NewNewAccount(AC("eosio"), AC("eosio.token"), ecc.MustNewPublicKey("EOS8ju96GnaKaYAs7b5EvvtwWqTVPepSCciDHvDCjiEhGb5joYtjk")),
		//system.NewNewAccount(AC("eosio"), AC("bilcproducer"), ecc.MustNewPublicKey("EOS8ju96GnaKaYAs7b5EvvtwWqTVPepSCciDHvDCjiEhGb5joYtjk")),
		//system.NewNewAccount(AC("eosio"), AC("cbillett"), ecc.MustNewPublicKey("EOS8ju96GnaKaYAs7b5EvvtwWqTVPepSCciDHvDCjiEhGb5joYtjk")),

		//bios
		//system.NewSetPriv(AC("eosio")),
		//system.NewSetPriv(AC("eosio.msig")),
		//system.NewSetPriv(AC("eosio.token")),
		//token.NewCreate(AC("eosio"), eos.NewEOSAsset(1000000000.0000), false, false, false),
		//token.NewIssue(AC("eosio"), eos.NewEOSAsset(1000000000.0000), ""),

		token.NewTransfer(eos.AccountName("eosio"), eos.AccountName("cbillett"), eos.NewEOSAsset(100000), ""),

	//system.NewRegProducer(
	//	AC("bilcproducer"),
	//	ecc.MustNewPublicKey("EOS8ju96GnaKaYAs7b5EvvtwWqTVPepSCciDHvDCjiEhGb5joYtjk"),
	//	system.EOSIOParameters{
	//		BasePerTransactionNetUsage:     100,
	//		BasePerTransactionCPUUsage:     500,
	//		BasePerActionCPUUsage:          1000,
	//		BaseSetcodeCPUUsage:            2097152, //# 2 * 1024 * 1024 // overbilling cpu usage for setcode to cover incidental
	//		PerSignatureCPUUsage:           100000,
	//		PerLockNetUsage:                32,
	//		ContextFreeDiscountCPUUsageNum: 20,
	//		ContextFreeDiscountCPUUsageDen: 100,
	//		MaxTransactionCPUUsage:         10485760,  //10 * 1024 * 1024
	//		MaxTransactionNetUsage:         102400,    // 100 * 1024
	//		MaxBlockCPUUsage:               104857600, // 100 * 1024 * 1024; at 500ms blocks and 20000instr trx, this enables ~10,000 TPS burst
	//		TargetBlockCPUUsagePct:         1000,      // 10%, 2 decimal places
	//		MaxBblockNetUsage:              1048576,
	//		TargetBlockNetUsagePct:         1000, // 10%, 2 decimal places
	//		MaxTransactionLifetime:         3600,
	//		MaxTransactionExecTime:         0, //unused??
	//		MaxAuthorityDepth:              6,
	//		MaxInlineDepth:                 4,
	//		MaxInlineActionSize:            4096,
	//		MaxGeneratedTransactionCount:   16,
	//		PercentOfMaxInflationRate:      10000, // percent, with 2 dec places.
	//		//MaxStorageSize: 10485760,  //FIXME
	//		StorageReserveRatio:            1000,  // ratio * 1000,
	//	},
	//),

	//system.NewVoteProducer(AC("cbillett"), AC(""), AC("bilcproducer")),
	//system.NewSetProds(1, []system.ProducerKey{
	//	{
	//			AC("eosio"),
	//			ecc.MustNewPublicKey("EOS8ju96GnaKaYAs7b5EvvtwWqTVPepSCciDHvDCjiEhGb5joYtjk"),
	//	},
	//}),

	)
	if err != nil {
		fmt.Println("ERROR calling NewAccount:", err)
	} else {
		fmt.Println("RESP:", actionResp)
	}

	//resp, err := api.GetCurrencyBalance(AC("eosio"), eos.EOSSymbol.Symbol, AC("eosio"))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(resp)

}
