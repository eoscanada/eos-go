package eos

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	mockserver "github.com/eoscanada/eos-go/testdata/mock_server"
	"github.com/stretchr/testify/assert"
)

var api *API

func TestAPIGetAccount(t *testing.T) {
	name := AccountName("teamgreymass")
	acc, err := api.GetAccount(context.Background(), name)
	assert.NoError(t, err)

	actualJSON, err := json.Marshal(acc)
	assert.NoError(t, err)

	expectedJSON := mockserver.OpenFile(".", "chain_get_account.json")

	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

func TestAPIGetInfo(t *testing.T) {
	info, err := api.GetInfo(context.Background())
	assert.NoError(t, err)

	actualJSON, err := json.Marshal(info)
	assert.NoError(t, err)

	expectedJSON := mockserver.OpenFile(".", "chain_get_info.json")

	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()

	os.Exit(code)
}

func setUp() {
	SetFloat64MarshalingTypeIntoString()

	mockserver.CreateAndActivateRestMockServer(".")

	api = New("http://localhost")

	// for working httpmock
	api.HttpClient = &http.Client{}
}

func tearDown() {
	mockserver.DeactivateMockServer()
}
