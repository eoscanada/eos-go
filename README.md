EOS.IO API library for Go
=========================

[点击查看中文版](./README-cn.md)

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

Binaries
--------

There is some binaries in `main` packages under `cmd/`, mainly around P2P communication.

Example
-------

### Reference

 * API
    * [Get Chain Information](./example_api_get_info_test.go)
    * [Transfer Token](./example_api_transfer_eos_test.go)
 * Decoding/Encoding
    * [Decode Table Row](./example_abi_decode_test.go)

### Running

The easiest way to see the actual output for a given example is to add a line
`// Output: any` at the very end of the test, looks like this for
`ExampleAPI_GetInfo` file ([examples_api_get_info.go](./examples_api_get_info.go)):

```
    if err != nil {
        panic(fmt.Errorf("json marshal response: %s", err))
    }

    fmt.Println(string(bytes))
    // Output: any
}
```

This tells `go test` that it can execute this test correctly. Then, simply
run only this example:

    go test -run ExampleAPI_GetInfo

Replacing `ExampleAPI_GetInfo` with the actual example name you want to try
out where line `// Output: any` was added.

This will run the example and compares the standard output with the `any` which
will fail. But it's ok an expected, so you can see the actual output
printed to your terminal.

**Note** Some examples will not succeed out of the box because it requires
some configuration. A good example being the `transfer` operation which
requires having the authorizations and balance necessary to perform the
transaction. It's quite possible to run them through a development environment
however.

#### Environment Variables

All examples uses by default the `https://mainnet.eos.dfuse.io` API endpoint for all
HTTP communication and `peering.mainnet.eoscanada.com` for P2P communication.
They can respectively be overridden by specifying environment variable
`EOS_GO_API_URL` and `EOS_GO_P2P_ENDPOINT` respectively.

Contributing
------------

Any contributions are welcome, use your standard GitHub-fu to pitch in and improve.


License
-------

MIT
