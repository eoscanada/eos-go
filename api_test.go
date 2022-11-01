package eos

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	mockserver "github.com/eoscanada/eos-go/testdata/mock_server"
	"github.com/stretchr/testify/assert"
)

var api *API

func TestAPIGetInfo(t *testing.T) {
	info, err := api.GetInfo(context.Background())
	if err != nil {
		panic(fmt.Errorf("get info: %w", err))
	}

	assert.NoError(t, err)
	assert.NotNil(t, info)
}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()

	os.Exit(code)
}

func setUp() {
	mockserver.CreateAndActivateRestMockServer()

	api = New("http://localhost")

	// for working httpmock
	api.HttpClient = &http.Client{}
}

func tearDown() {
	mockserver.DeactivateMockServer()
}
