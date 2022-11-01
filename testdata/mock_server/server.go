package mockServer

import (
	"net/http"
	"os"
	"strings"

	"github.com/jarcoal/httpmock"
)

const (
	GET  = "GET"
	POST = "POST"
)

func CreateAndActivateRestMockServer() {
	httpmock.Activate()

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_account",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_account.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_block",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_block.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_block_info",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_block_info.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_info",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_info.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_block_header_state",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_block_header_state.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_abi",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_abi.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_currency_balance",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_currency_balance.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_currency_stats",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_currency_stats.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_required_keys",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_required_keys.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_producers",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_producers.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_raw_code_and_abi",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_raw_code_and_abi.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_table_by_scope",
		httpmock.NewStringResponder(http.StatusOK, openFile("chain_get_table_by_scope.json")),
	)
}

func openFile(filename string) string {
	path := []string{"./testdata/mock_server", filename}
	body, err := os.ReadFile(strings.Join(path, "/"))
	if err != nil {
		panic(err)
	}

	return string(body)
}

func DeactivateMockServer() {
	httpmock.DeactivateAndReset()
}
