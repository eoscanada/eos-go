## 用 Go 语言与 EOS.IO 交互的 API 库

[![GoDoc](https://godoc.org/github.com/eoscanada/eos-go?status.svg)](https://godoc.org/github.com/eoscanada/eos-go)

该库提供对数据架构（二进制打包和JSON接口）的简单访问，
以及对远程或本地运行的EOS.IO RPC服务器的API调用。
它提供钱包功能（KeyBag），或者可以通过 `keosd` 钱包签署交易。
它还明白端口9876上的P2P协议。

截至6月的发布之前，这个库不断的在变化。 先不要期望稳定性，
因为我们要追着主网 `eosio` 代码库的脚步，而它的变化又那么快。

该库主网启动编排工具是 `eosio` 的基础，网址：
https://github.com/eoscanada/eos-bios

### 基本用法

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

### 例子

#### 参考

 * API
    * [获取链信息](./example_api_get_info_test.go)
    * [转账代币](./example_api_transfer_eos_test.go)
 * 解码/编码
    * [解码表行](./example_abi_decode_test.go)

### 二进制文件

`cmd/` 下的 `main` 包中有一些二进制文件，主要围绕 P2P 通信。

### 召集开源贡献者

我们欢迎所有的开源贡献，直接用 GitHub-fu来提议、帮我们改进吧。

### 证书

MIT
