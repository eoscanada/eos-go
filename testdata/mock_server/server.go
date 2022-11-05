package mockServer

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/jarcoal/httpmock"
)

const (
	GET  = "GET"
	POST = "POST"
)

func CreateAndActivateRestMockServer(rootDir string) {
	httpmock.Activate()

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_account",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_account.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_block",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_block.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_block_info",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_block_info.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_info",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_info.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_block_header_state",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_block_header_state.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_abi",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_abi.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_currency_balance",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_currency_balance.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_currency_stats",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_currency_stats.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_required_keys",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_required_keys.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_producers",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_producers.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_raw_code_and_abi",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_raw_code_and_abi.json")),
	)

	httpmock.RegisterResponder(
		POST, "http://localhost/v1/chain/get_table_by_scope",
		httpmock.NewStringResponder(http.StatusOK, OpenFile(rootDir, "chain_get_table_by_scope.json")),
	)
}

func OpenFile(rootdir, filename string) string {
	body, err := os.ReadFile(filepath.Join(rootdir, "testdata", "mock_server", filename))
	if err != nil {
		panic(err)
	}

	return string(body)
}

func DeactivateMockServer() {
	httpmock.DeactivateAndReset()
}
