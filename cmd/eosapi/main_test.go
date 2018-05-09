package main

import (
	"net/http"
	"testing"
	"time"

	"bytes"

	"log"

	"fmt"

	"encoding/hex"

	"github.com/eoscanada/eos-go"
	"github.com/gin-gonic/gin/json"
	"github.com/stretchr/testify/assert"
)

func newAPI() (api *eos.API) {

	api = eos.New("http://localhost:8888", bytes.Repeat([]byte{0}, 32))
	//api = eos.New("http://stage0.eoscanada.com", bytes.Repeat([]byte{0}, 32))

	tr := &http.Transport{}
	api.HttpClient = &http.Client{Transport: tr}
	keyBag := eos.NewKeyBag()

	for _, key := range []string{
		"5HryYjdRzBtQKzM1H7L1Y4yokBMAoUYjcYpMvhQv1hzKhrKdfWp",
		"5KALFt6C4nntzAsT1cDiiPgw1cdDJPzutadSCFLygqGvKBjvNqP",
	} {
		if err := keyBag.Add(key); err != nil {
			log.Fatalln("Couldn't load private key:", err)
		}
	}

	api.SetSigner(keyBag)

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

	out, err := api.GetTableRows(eos.GetTableRowsRequest{
		Scope:    "currency",
		Code:     "currency",
		Table:    "account",
		TableKey: "currency",
		JSON:     false,
	})

	assert.NoError(t, err)
	assert.NotNil(t, out.Rows)

}

func TestGetTransactions(t *testing.T) {
	api := newAPI()
	api.Debug = true
	out, err := api.GetTransactions(eos.AccountName("eosio.disco"))
	assert.NoError(t, err)
	data, err := json.Marshal(out)
	fmt.Println(string(data))
}

func TestGetTransactionsForValidateChain(t *testing.T) {
	api := newAPI()
	api.Debug = true
	accounts := []eos.AccountName{
		//eos.AccountName("eosio"),
		eos.AccountName("eosio.msig"),
		eos.AccountName("eosio.disco"),
		eos.AccountName("eosio.token"),
	}
	actions := make(map[string]*eos.Action, 0)
	for _, account := range accounts {

		out, err := api.GetTransactions(account)
		assert.NoError(t, err)
		for _, tx := range out.Transactions {
			for _, action := range tx.Transaction.Transaction.Actions {

				key := hex.EncodeToString(action.HexData)
				fmt.Println("key : ", key)
				data, err := json.Marshal(action)
				assert.NoError(t, err)
				fmt.Println("Data  : ", string(data))
				if collision, ok := actions[key]; ok {
					cdata, err := json.Marshal(collision)
					assert.NoError(t, err)
					fmt.Println("CData : ", string(cdata))
					fmt.Println("Found a colision")
				}
				actions[key] = action
			}
		}
	}

	fmt.Printf("Found [%d] actions\n", len(actions))
}

func TestGetRequiredKeys(t *testing.T) {
	api := newAPI()
	tomorrow := time.Now().AddDate(0, 0, 1)
	keybag := eos.NewKeyBag()
	api.SetSigner(keybag)
	out, err := api.GetRequiredKeys(&eos.Transaction{
		// RefBlockNum:    "1",
		// RefBlockPrefix: "",
		Expiration: eos.JSONTime{tomorrow},
		// Scope:          []string{},
		// Actions: []eos.Action{
		// 	{
		// 		Account: eos.AccountName("currency"),
		// 		Name:    "currency",
		// 		// Authorization: []string{},
		// 		Data: "",

		// 		Type: "dawn-2",
		// 		Code: "dawn-2",
		// 		// Recipients: []string{"currency"},
		// 	},
		// },
		// Signatures:     []string{},
		// Authorizations: []string{},
	})
	assert.NoError(t, err)
	assert.Equal(t, "mama", out.RequiredKeys)
}

func TestGetCurrencyBalance(t *testing.T) {
	api := newAPI()

	out, err := api.GetCurrencyBalance("currency", "CUR", "currency")
	assert.NoError(t, err)
	assert.NotZero(t, len(out))
}
