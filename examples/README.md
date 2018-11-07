## `eos-go` Examples

This folder contains a bunch of examples to perform some tasks with the
library.

### Reference

 * [Get chain information](./get_info/main.go)
 * [Transfer token](./transfer/main.go)

### Running

Examples can be run simply by doing:

    go run examples/get_info/main.go

**Note** Some examples will not succeed out of the box because it requires
some configuration. A good example being the `transfer` operation which
requires having the authorizations and balance necessary to perform the
transaction. It's quite possible to run them through a development environment
however.

#### Environment Variables

All examples uses the `API_URL` environment to talk to the API endpoint, when
not provided, they use by default
