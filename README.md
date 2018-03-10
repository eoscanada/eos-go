EOS.IO API library for Go
=========================

[![GoDoc](https://godoc.org/github.com/eosioca/eosapi?status.svg)](https://godoc.org/github.com/eosioca/eosapi)

This library aims to provide simple access to data structures and API
calls to an EOS.IO RPC server, running remotely or locally.

Basic usage
-----------

```go
api := eosapi.New("http://testnet1.eos.io")

infoResp, _ := api.GetInfo()
accountResp, _ := api.GetAccount("initn")
fmt.Println("Permission for initn:", accountResp.Permissions[0].RequiredAuth.Keys)
```

Contributing
------------

Any contributions are welcome, use your standard GitHub-fu to pitch in and improve.


License
-------

MIT
