EOS.IO API library for Go
=========================

[![GoDoc](https://godoc.org/github.com/eoscanada/eos-go?status.svg)](https://godoc.org/github.com/eoscanada/eos-go)

This library provides simple access to data structures (binary packing
and JSON interface) and API calls to an EOS.IO RPC server, running
remotely or locally.  It provides wallet functionalities (KeyBag), or
can sign transaction through the `keosd` wallet. It also knows about
the P2P protocol on port 9876.

As of before the June launch, this library is pretty much in
flux. Don't expect stability, as we're moving alongside the main
`eosio` codebase, which changes very fast.

This library is the basis for the `eos-bios` launch orchestrator tool
at https://github.com/eoscanada/eos-bios


Basic usage
-----------

```go
api := eos.New("http://testnet1.eos.io")

infoResp, _ := api.GetInfo()
accountResp, _ := api.GetAccount("initn")
fmt.Println("Permission for initn:", accountResp.Permissions[0].RequiredAuth.Keys)
```

`eosio.system` and `eosio.token` contract _Actions_ are respectively in:
* https://github.com/eoscanada/eos-go/tree/master/system ([godocs](https://godoc.org/github.com/eoscanada/eos-go/system))
* https://github.com/eoscanada/eos-go/tree/master/token ([godocs](https://godoc.org/github.com/eoscanada/eos-go/token))

Example
-------

See example usages of the library:

* https://github.com/eoscanada/eos-bios/blob/master/bios.go
* https://github.com/eoscanada/eos-bios/blob/master/ops.go
* Some other `main` packages under `cmd/`.


Contributing
------------

Any contributions are welcome, use your standard GitHub-fu to pitch in and improve.


License
-------

MIT



----------------------

Changes to dawn4:
* sig_digest always adds something even with empty context free actions.
* PUB, PVT, SIG
* implement `delegatebw` NewDelegateBandwidth(), etc..
  (storage_stake)(storage_bytes) -> (ram_bytes)

Unanswered questions:
* what's the "location" field on a "producer_info" (eosio.system.abi)
* the `setglimits` and `setalimits` in the eosio.bios .. do we need to call that
  to setup the chain ? what does bootseq_.. say in the eosio repo ?
* what do `transaction_extensions` mean in a transaction ? any uses ? special cases ?
* all the multisig things.. we should implement
* usage is billed, means the "scope" of the storage on those multi_index tables is "contract" and "billed account", right ?! :) that's pretty cool :)

* does --stake-net and friends create a separate action ?
* cleos create newaccount --stake-net "1.0000 EOS" --stake-cpu "1.0000 EOS" --buy-ram-bytes 111 eoosio acct pubkey pubkey -p eosio
* check `eosio` tables now: producers, global, voters, userres, totalband, delband, refunds, and `msig` table `proposal`

* sync up `get_account` return struct.. way richer now.
* sync up most /v1/chain outputs..
