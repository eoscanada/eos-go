package eos

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/eoscanada/eos-go/ecc"
)

type API struct {
	HttpClient              *http.Client
	BaseURL                 string
	Signer                  Signer
	Debug                   bool
	Compress                CompressionType
	DefaultMaxCPUUsageMS    uint8
	DefaultMaxNetUsageWords uint32 // in 8-bytes words

	lastGetInfo      *InfoResp
	lastGetInfoStamp time.Time
	lastGetInfoLock  sync.Mutex

	customGetRequiredKeys func(tx *Transaction) ([]ecc.PublicKey, error)
}

func New(baseURL string) *API {
	api := &API{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     true, // default behavior, because of `nodeos`'s lack of support for Keep alives.
			},
		},
		BaseURL:  baseURL,
		Compress: CompressionZlib,
	}

	return api
}

// FixKeepAlives tests the remote server for keepalive support (the
// main `nodeos` software doesn't in the version from March 22nd
// 2018).  Some endpoints front their node with a keep-alive
// supporting web server.  Adjust the `KeepAlive` support of the
// client accordingly.
func (api *API) FixKeepAlives() bool {
	// Yeah, to provoke a keep alive, you need to query twice.
	for i := 0; i < 5; i++ {
		_, err := api.GetInfo()
		if api.Debug {
			log.Println("err", err)
		}
		if err == io.EOF {
			if tr, ok := api.HttpClient.Transport.(*http.Transport); ok {
				tr.DisableKeepAlives = true
				return true
			}
		}
		_, err = api.GetNetConnections()
		if api.Debug {
			log.Println("err", err)
		}
		if err == io.EOF {
			if tr, ok := api.HttpClient.Transport.(*http.Transport); ok {
				tr.DisableKeepAlives = true
				return true
			}
		}
	}
	return false
}

func (api *API) EnableKeepAlives() bool {
	if tr, ok := api.HttpClient.Transport.(*http.Transport); ok {
		tr.DisableKeepAlives = false
		return true
	}
	return false
}

func (api *API) SetCustomGetRequiredKeys(f func(tx *Transaction) ([]ecc.PublicKey, error)) {
	api.customGetRequiredKeys = f
}

func (api *API) SetSigner(s Signer) {
	api.Signer = s
}

// ProducerPause will pause block production on a nodeos with
// `producer_api` plugin loaded.
func (api *API) ProducerPause() error {
	return api.call("producer", "pause", nil, nil)
}

// ProducerResume will resume block production on a nodeos with
// `producer_api` plugin loaded. Obviously, this needs to be a
// producing node on the producers schedule for it to do anything.
func (api *API) ProducerResume() error {
	return api.call("producer", "resume", nil, nil)
}

// IsProducerPaused queries the blockchain for the pause statement of
// block production.
func (api *API) IsProducerPaused() (out bool, err error) {
	err = api.call("producer", "paused", nil, &out)
	return
}

func (api *API) GetAccount(name AccountName) (out *AccountResp, err error) {
	err = api.call("chain", "get_account", M{"account_name": name}, &out)
	return
}

func (api *API) GetCode(account AccountName) (out *GetCodeResp, err error) {
	err = api.call("chain", "get_code", M{"account_name": account, "code_as_wasm": true}, &out)
	return
}

func (api *API) GetABI(account AccountName) (out *GetABIResp, err error) {
	err = api.call("chain", "get_abi", M{"account_name": account}, &out)
	return
}

// WalletImportKey loads a new WIF-encoded key into the wallet.
func (api *API) WalletImportKey(walletName, wifPrivKey string) (err error) {
	return api.call("wallet", "import_key", []string{walletName, wifPrivKey}, nil)
}

func (api *API) WalletPublicKeys() (out []ecc.PublicKey, err error) {
	var textKeys []string
	err = api.call("wallet", "get_public_keys", nil, &textKeys)
	if err != nil {
		return nil, err
	}

	for _, k := range textKeys {
		newKey, err := ecc.NewPublicKey(k)
		if err != nil {
			return nil, err
		}

		out = append(out, newKey)
	}
	return
}

func (api *API) ListKeys() (out []*ecc.PrivateKey, err error) {
	var textKeys []string
	err = api.call("wallet", "list_keys", nil, &textKeys)
	if err != nil {
		return nil, err
	}

	for _, k := range textKeys {
		newKey, err := ecc.NewPrivateKey(k)
		if err != nil {
			return nil, err
		}

		out = append(out, newKey)
	}
	return
}

func (api *API) WalletSignTransaction(tx *SignedTransaction, chainID []byte, pubKeys ...ecc.PublicKey) (out *WalletSignTransactionResp, err error) {
	var textKeys []string
	for _, key := range pubKeys {
		textKeys = append(textKeys, key.String())
	}

	err = api.call("wallet", "sign_transaction", []interface{}{
		tx,
		textKeys,
		hex.EncodeToString(chainID),
	}, &out)
	return
}

func (api *API) SignPushActions(a ...*Action) (out *PushTransactionFullResp, err error) {
	return api.SignPushTransaction(&Transaction{Actions: a}, nil)
}

func (api *API) SignPushActionsWithOpts(opts TxOptions, a ...*Action) (out *PushTransactionFullResp, err error) {
	return api.SignPushTransaction(&Transaction{Actions: a}, &opts)
}

func (api *API) SignPushTransaction(tx *Transaction, opts *TxOptions) (out *PushTransactionFullResp, err error) {
	if api.Signer == nil {
		return nil, fmt.Errorf("no Signer configured")
	}

	if opts == nil {
		opts = &TxOptions{}
	}

	_, err = tx.Fill(api)
	if err != nil {
		return nil, err
	}

	var requiredKeys []ecc.PublicKey
	if api.customGetRequiredKeys != nil {
		requiredKeys, err = api.customGetRequiredKeys(tx)
		if err != nil {
			return nil, fmt.Errorf("custom_get_required_keys: %s", err)
		}
	} else {
		resp, err := api.GetRequiredKeys(tx)
		if err != nil {
			return nil, fmt.Errorf("get_required_keys: %s", err)
		}
		requiredKeys = resp.RequiredKeys
	}

	stx := NewSignedTransaction(tx)

	stx.estimateResources(*opts, api.DefaultMaxCPUUsageMS, api.DefaultMaxNetUsageWords)

	signedTx, err := api.Signer.Sign(stx, api.lastGetInfo.ChainID, requiredKeys...)
	if err != nil {
		return nil, fmt.Errorf("signing through wallet: %s", err)
	}

	packed, err := signedTx.Pack(*opts)
	if err != nil {
		return nil, err
	}

	return api.PushSignedTransaction(packed)
}

func (api *API) PushSignedTransaction(tx *PackedTransaction) (out *PushTransactionFullResp, err error) {
	err = api.call("chain", "push_transaction", tx, &out)
	return
}

func (api *API) GetInfo() (out *InfoResp, err error) {
	err = api.call("chain", "get_info", nil, &out)
	return
}

func (api *API) GetNetConnections() (out []*NetConnectionsResp, err error) {
	err = api.call("net", "connections", nil, &out)
	return
}

func (api *API) NetConnect(host string) (out NetConnectResp, err error) {
	err = api.call("net", "connect", host, &out)
	return
}

func (api *API) NetDisconnect(host string) (out NetDisconnectResp, err error) {
	err = api.call("net", "disconnect", host, &out)
	return
}

func (api *API) GetNetStatus(host string) (out *NetStatusResp, err error) {
	err = api.call("net", "status", M{"host": host}, &out)
	return
}

func (api *API) GetBlockByID(id string) (out *BlockResp, err error) {
	err = api.call("chain", "get_block", M{"block_num_or_id": id}, &out)
	return
}

func (api *API) GetProducers() (out *ProducersResp, err error) {
	/*
		+FC_REFLECT( eosio::chain_apis::read_only::get_producers_params, (json)(lower_bound)(limit) )
		+FC_REFLECT( eosio::chain_apis::read_only::get_producers_result, (rows)(total_producer_vote_weight)(more) ); */
	err = api.call("chain", "get_producers", nil, &out)
	return
}

func (api *API) GetBlockByNum(num uint32) (out *BlockResp, err error) {
	err = api.call("chain", "get_block", M{"block_num_or_id": fmt.Sprintf("%d", num)}, &out)
	//err = api.call("chain", "get_block", M{"block_num_or_id": num}, &out)
	return
}

func (api *API) GetBlockByNumOrID(query string) (out *SignedBlock, err error) {
	err = api.call("chain", "get_block", M{"block_num_or_id": query}, &out)
	return
}

func (api *API) GetTransaction(id string) (out *TransactionResp, err error) {
	err = api.call("account_history", "get_transaction", M{"transaction_id": id}, &out)
	return
}

func (api *API) GetTransactions(name AccountName) (out *TransactionsResp, err error) {
	err = api.call("account_history", "get_transactions", M{"account_name": name}, &out)
	return
}

func (api *API) GetTableRows(params GetTableRowsRequest) (out *GetTableRowsResp, err error) {
	err = api.call("chain", "get_table_rows", params, &out)
	return
}

func (api *API) GetRequiredKeys(tx *Transaction) (out *GetRequiredKeysResp, err error) {
	keys, err := api.Signer.AvailableKeys()
	if err != nil {
		return nil, err
	}

	err = api.call("chain", "get_required_keys", M{"transaction": tx, "available_keys": keys}, &out)
	return
}

func (api *API) GetCurrencyBalance(account AccountName, symbol string, code AccountName) (out []Asset, err error) {
	err = api.call("chain", "get_currency_balance", M{"account": account, "symbol": symbol, "code": code}, &out)
	return
}

// See more here: libraries/chain/contracts/abi_serializer.cpp:58...

func (api *API) call(baseAPI string, endpoint string, body interface{}, out interface{}) error {
	jsonBody, err := enc(body)
	if err != nil {
		return err
	}

	targetURL := fmt.Sprintf("%s/v1/%s/%s", api.BaseURL, baseAPI, endpoint)
	req, err := http.NewRequest("POST", targetURL, jsonBody)
	if err != nil {
		return fmt.Errorf("NewRequest: %s", err)
	}

	if api.Debug {
		// Useful when debugging API calls
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("-------------------------------")
		fmt.Println(string(requestDump))
		fmt.Println("")
	}

	resp, err := api.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %s", req.URL.String(), err)
	}
	defer resp.Body.Close()

	var cnt bytes.Buffer
	_, err = io.Copy(&cnt, resp.Body)
	if err != nil {
		return fmt.Errorf("Copy: %s", err)
	}

	if resp.StatusCode == 404 {
		return ErrNotFound
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("%s: status code=%d, body=%s", req.URL.String(), resp.StatusCode, cnt.String())
	}

	if api.Debug {
		fmt.Println("RESPONSE:")
		fmt.Println(cnt.String())
		fmt.Println("")
	}

	if err := json.Unmarshal(cnt.Bytes(), &out); err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}

	return nil
}

var ErrNotFound = errors.New("resource not found")

type M map[string]interface{}

func enc(v interface{}) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}

	cnt, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	//fmt.Println("BODY", string(cnt))

	return bytes.NewReader(cnt), nil
}
