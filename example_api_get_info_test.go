package eos_test

import (
	"context"
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
)

func ExampleAPI_GetInfo() {
	api := eos.New(getAPIURL())

	info, err := api.GetInfo(context.Background())
	if err != nil {
		panic(fmt.Errorf("get info: %w", err))
	}

	bytes, err := json.Marshal(info)
	if err != nil {
		panic(fmt.Errorf("json marshal response: %w", err))
	}

	fmt.Println(string(bytes))
}
