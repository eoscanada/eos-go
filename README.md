## EOS.IO API library for Go

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


### Basic usage

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	cli "github.com/streamingfast/cli"
)

func main() {
	api := eos.New("https://api.eosn.io")
	ctx := context.Background()

	infoResp, err := api.GetInfo(ctx)
	cli.NoError(err, "unable to get chain info")

	fmt.Println("Chain Info", toJson(infoResp))

	accountResp, _ := api.GetAccount(ctx, "eosio")
	fmt.Println("Account Info", toJson(accountResp))
}

func toJson(v interface{}) string {
	out, err := json.MarshalIndent(v, "", "  ")
	cli.NoError(err, "unable to marshal json")

	return string(out)
}
```

### Examples

#### Reference

 * API
	* [Get Account](./example_api_get_account_test.go)
	* [Get Chain Information](./example_api_get_info_test.go)
	* [Get Producers](./example_api_get_producers_test.go)
	* [Transfer Token](./example_api_transfer_eos_test.go)
 * Decoding/Encoding
	* [Decode Table Row](./example_abi_decode_test.go)
 * Transaction
	* [Transaction Sign & Pack](./example_trx_pack_test.go)
	* [Transaction Unpack](./example_trx_unpack_test.go)

#### Running

The easiest way to see the actual output for a given example is to add a line
`// Output: any` at the very end of the test, looks like this for
`ExampleAPI_GetInfo` file ([examples_api_get_info.go](./example_api_get_info_test.go)):

```
	...

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

### Binaries

There is some binaries in `main` packages under `cmd/`, mainly around P2P communication.

### Environment Variables

All examples uses by default the `https://api.eosn.io` API endpoint for all
HTTP communication and `peering.eosn.io` for P2P communication.
They can respectively be overridden by specifying environment variable
`EOS_GO_API_URL` and `EOS_GO_P2P_ENDPOINT` respectively.

### Tests

Some of our tests renders dates in the timezone of the OS. As such, if you have a bunch of
failures around dates and times, it's probably because your timezone is not aligned with
those in the tests.

Run the tests with this to be in the same timezone as the expected one in golden files:

```bash
TZ=UTC go test ./...
```

### Release

We are using [Goreleaser](https://goreleaser.com/) to perform releases. Install the `goreleaser` binary ([instructions](https://goreleaser.com/install/))
and follow these steps:

- Dry run release process first with `goreleaser release --skip-publish --skip-validate --rm-dist`
- Publish **draft** release with `goreleaser release --rm-dist`
- Open GitHub's release and check that the release is all good
- Once everything is good, publish release, this will now also push the tag.

### Contributing

Any contributions are welcome, use your standard GitHub-fu to pitch in and improve.

### License

MIT
