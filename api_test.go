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

func TestAPIGetBlockByNum(t *testing.T) {
	block, err := api.GetBlockByNum(context.Background(), 273283700)
	assert.NoError(t, err)

	assert.Equal(t, block.ID.String(), "1049fa74670c8d45c83cfd6b54683edb186b5205bf84c66141afe55765499f7c")
	assert.Equal(t, block.RefBlockPrefix, uint32(1811758280))
	assert.Equal(t, block.Timestamp.Format("2006-01-02T15:04:05"), "2022-10-15T07:13:42")
	assert.Equal(t, block.Producer.String(), "eosiosg11111")
	assert.Equal(t, block.Confirmed, uint16(240))
	assert.Equal(t, block.Previous.String(), "1049fa737265f6c2d46428c2ee89cd5be29dc2b9a0615afc6af4c8df6ffa7845")
	assert.Equal(t, block.TransactionMRoot.String(), "f7cadbcc33efad01eacc153d2013ce9e50fe4f0d664df657a354a3d06ac569ab")
	assert.Equal(t, block.ActionMRoot.String(), "0378d11e50e8ee4c8673a06ece81933d672ff0578aec743aad4f6962309dca58")
	assert.Equal(t, block.ScheduleVersion, uint32(2043))
	assert.Nil(t, block.NewProducersV1)
	assert.Equal(t, block.ProducerSignature.String(), "SIG_K1_KbbKB47LguWUfhYfqcZPNTR6Nd8hgLV4CfF1GSxLf7TBqwaFDbbXdMvDDQKhp5186HSWz1MgmNt2qPHVWdY6KAkZaDbEab")

	assert.Equal(t, len(block.Transactions), 1)

	tx := block.Transactions[0]

	assert.Equal(t, tx.Status.String(), "executed")
	assert.Equal(t, tx.CPUUsageMicroSeconds, uint32(389))
	assert.Equal(t, tx.NetUsageWords, Varuint32(33))
	assert.Equal(t, tx.Transaction.ID.String(), "3b842c3b6eb260028a51bc9c4b1cf9587b393a0607b45dafd7c5279c200c3e24")
	assert.Equal(t, len(tx.Transaction.Packed.Signatures), 2)
	assert.Equal(t, tx.Transaction.Packed.Signatures[0].String(), "SIG_K1_KjUMXRq5vgxs9xenpjCR1PBP5vNQNndVD5HyXRtqjrQL4h7NRaS6iVhRXtdst6J4fYxnhbbfnbJsXoWiPugoVU8DZBoG2o")
	assert.Equal(t, tx.Transaction.Packed.Compression.String(), "none")
	assert.Equal(t, tx.Transaction.Packed.ContextFreeData, []string{})
	assert.Equal(t, tx.Transaction.Packed.PackedTransaction.String(),
		"ba5d4a6360fae072f0530000000001a026a59a4d8331550080cae6aa4addd402a02bd21551cda6c100000000a8ed3232e07ba59a4d83315500000000a8ed3232880140aeda34d25cfd450000c8d7645cbb920000000000000000896f0000001f2819b684a6cff1e72d17a0e050a3349797dea09cbaef132f241790f2c73d49b9633cb6a595b7cd16378ded6323dbc131456cdd7ced6b287ad93a24bc1945771a10000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, tx.Transaction.Packed.Transaction.Expiration.Format("2006-01-02T15:04:05"), "2022-10-15T07:14:02")
	assert.Equal(t, tx.Transaction.Packed.Transaction.RefBlockNum, uint16(64096))
	assert.Equal(t, tx.Transaction.Packed.Transaction.RefBlockPrefix, uint32(1408266976))
	assert.Equal(t, tx.Transaction.Packed.Transaction.MaxNetUsageWords, Varuint32(0))
	assert.Equal(t, tx.Transaction.Packed.Transaction.MaxCPUUsageMS, uint8(0))
	assert.Equal(t, tx.Transaction.Packed.Transaction.DelaySec, Varuint32(0))
	assert.Equal(t, tx.Transaction.Packed.Transaction.ContextFreeActions, []*Action{})

	assert.Equal(t, len(tx.Transaction.Packed.Transaction.Actions), 1)

	action := tx.Transaction.Packed.Transaction.Actions[0]
	assert.Equal(t, action.Account.String(), "eossanguoone")
	assert.Equal(t, action.Name.String(), "unioperate")
	assert.Equal(t, action.Authorization[0].Actor.String(), "sanguocpucpu")
	assert.Equal(t, action.Authorization[0].Permission.String(), "active")

	// That's the main reason why I assert strictly, but no fancy way to check..
	assert.NotNil(t, action.Data)
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
