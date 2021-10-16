package eos_test

import (
	"context"
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
)

func ExampleAPI_GetAccount() {
	api := eos.New(getAPIURL())

	account := eos.AccountName("eos.rex")
	info, err := api.GetAccount(context.Background(), account)
	if err != nil {
		if err == eos.ErrNotFound {
			fmt.Printf("unknown account: %s", account)
			return
		}

		panic(fmt.Errorf("get account: %w", err))
	}

	bytes, err := json.Marshal(info)
	if err != nil {
		panic(fmt.Errorf("json marshal response: %w", err))
	}

	fmt.Println(string(bytes))
}
