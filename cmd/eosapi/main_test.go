package main

import (
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
	out, err := api.GetAccount("currency")
	assert.NoError(t, err)
	assert.NotNil(t, out.AccountName)
}

func TestGetCode(t *testing.T) {
	api := newAPI()
	out, err := api.GetCode("currency")
	assert.NoError(t, err)
	assert.Equal(t, "currency", out.AccountName)
}

func TestGetInfo(t *testing.T) {
	api := newAPI()
	out, err := api.GetInfo()
	assert.NoError(t, err)
	assert.Equal(t, "eosio", string(out.HeadBlockProducer))

}

func TestGetBlockByID(t *testing.T) {
	api := newAPI()
	blockID := "000244e5696fd9efadd76f1722ae683c9ea48a042392f5d3b7705b22ff5a11f9"
	out, err := api.GetBlockByID(blockID)
	assert.NoError(t, err)
	assert.Equal(t, blockID, out.ID)
}

func TestGetBlockByNum(t *testing.T) {
	api := newAPI()
	blockNum := uint64(1)
	out, err := api.GetBlockByNum(blockNum)
	assert.NoError(t, err)
	assert.Equal(t, blockNum, out.BlockNum)
}

func TestGetTableRows(t *testing.T) {
	api := newAPI()

	out, err := api.GetTableRows(eosapi.GetTableRowsRequest{
		Scope:    "currency",
		Code:     "currency",
		Table:    "account",
		TableKey: "currency",
		JSON:     true,
	})

	assert.NoError(t, err)
	assert.NotNil(t, out.Rows)

}

// func (api *EOSAPI) GetRequiredKeys(tx Transaction, availableKeys ...PublicKey) (out *GetRequiredKeysResp, err error) {
// 	err = api.call("chain", "get_required_keys", M{"transaction": tx, "available_keys": availableKeys}, &out)
// 	return
// }
// func TestGetRequiredKeys(t *testing.T) {
// 	api := newAPI()
// 	tomorrow := time.Now().AddDate(0, 0, 1)
// 	out, err := api.GetRequiredKeys(eosapi.Transaction{
// 		RefBlockNum:    "1",
// 		RefBlockPrefix: "",
// 		Expiration:     eosapi.JSONTime(tomorrow),
// 		Scope:          []string{},
// 		Actions: []eosapi.Action{
// 			Account:       eosapi.AccountName("currency"),
// 			Name:          "currency",
// 			Authorization: []string{},
// 			Data:          "",

// 			Type:       "dawn-2",
// 			Code:       "dawn-2",
// 			Recipients: []string{"currency"},
// 		},
// 		Signatures:     []string{},
// 		Authorizations: []string{},
// 	}, "EOS")
// 	assert.NoError(t, err)
// }

func TestGetCurrencyBalance(t *testing.T) {
	api := newAPI()

	out, err := api.GetCurrencyBalance("currency", "CUR", "currency")
	assert.NoError(t, err)
	assert.NotZero(t, len(out))
}
