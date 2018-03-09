package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/eosioca/eosapi"
	"github.com/stretchr/testify/assert"
)

func newAPI() (api *eosapi.EOSAPI) {
	api = eosapi.New("http://testnet-dawn3.eosio.ca")
	tr := &http.Transport{}
	api.HttpClient = &http.Client{Transport: tr}

	return
}

func TestGetAccount(t *testing.T) {
	api := newAPI()
	out, err := api.GetAccount("initm")
	assert.NoError(t, err)
	fmt.Println("Account initm", out)
}

func TestGetInfo(t *testing.T) {
	api := newAPI()
	out, err := api.GetInfo()
	assert.NoError(t, err)
	assert.Equal(t, "eosio", string(out.HeadBlockProducer))

}

func TestGetTableRows(t *testing.T) {
	api := newAPI()

	out, err := api.GetTableRows(eosapi.GetTableRowsRequest{
		Scope:    "currency",
		Code:     "currency",
		Table:    "account",
		TableKey: "currency",
		JSON:     false,
	})
	cnt, err := json.MarshalIndent(out, "", "  ")

	assert.NoError(t, err)
	assert.Equal(t, "", string(cnt))
	fmt.Println("GetBlockNum", out)

}

// func TestGetCurrencyBalance(t *testing.T) {
// 	api := newAPI()

// 	out, err := api.GetCurrencyBalance("initm", "EOS", "eosio")
// 	assert.NoError(t, err)

// 	cnt, err := json.MarshalIndent(out, "", "  ")
// 	assert.Equal(t, "[\n  \"1000004.0031 EOS\"\n]", string(cnt))
// }
