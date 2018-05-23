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
	ChainID                 []byte
	Signer                  Signer
	Debug                   bool
	Compress                CompressionType
	DefaultMaxCPUUsageMS    uint8
	DefaultMaxNetUsageWords uint32 // in 8-bytes words

	lastGetInfo      *InfoResp
	lastGetInfoStamp time.Time
	lastGetInfoLock  sync.Mutex
}

func New(baseURL string, chainID []byte) *API {
	if len(chainID) != 32 {
		panic("chainID must be 32 bytes")
	}

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
		ChainID:  chainID,
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

func (api *API) SetSigner(s Signer) {
	api.Signer = s
}

// Chain APIs
// Wallet APIs

// List here: https://github.com/Netherdrake/py-eos-api/blob/master/eosapi/client.py

// const string push_txn_func = chain_func_base + "/push_transaction";
// const string push_txns_func = chain_func_base + "/push_transactions";
// const string json_to_bin_func = chain_func_base + "/abi_json_to_bin";
// const string get_block_func = chain_func_base + "/get_block";
// const string get_account_func = chain_func_base + "/get_account";
// const string get_table_func = chain_func_base + "/get_table_rows";
// const string get_code_func = chain_func_base + "/get_code";
// const string get_currency_balance_func = chain_func_base + "/get_currency_balance";
// const string get_currency_stats_func = chain_func_base + "/get_currency_stats";
// const string get_required_keys = chain_func_base + "/get_required_keys";

// const string account_history_func_base = "/v1/account_history";
// const string get_transaction_func = account_history_func_base + "/get_transaction";
// const string get_transactions_func = account_history_func_base + "/get_transactions";
// const string get_key_accounts_func = account_history_func_base + "/get_key_accounts";
// const string get_controlled_accounts_func = account_history_func_base + "/get_controlled_accounts";

// const string net_func_base = "/v1/net";
// const string net_connect = net_func_base + "/connect";
// const string net_disconnect = net_func_base + "/disconnect";
// const string net_status = net_func_base + "/status";
// const string net_connections = net_func_base + "/connections";

// const string wallet_func_base = "/v1/wallet";
// const string wallet_create = wallet_func_base + "/create";
// const string wallet_open = wallet_func_base + "/open";
// const string wallet_list = wallet_func_base + "/list_wallets";
// const string wallet_list_keys = wallet_func_base + "/list_keys";
// const string wallet_public_keys = wallet_func_base + "/get_public_keys";
// const string wallet_lock = wallet_func_base + "/lock";
// const string wallet_lock_all = wallet_func_base + "/lock_all";
// const string wallet_unlock = wallet_func_base + "/unlock";
// const string wallet_import_key = wallet_func_base + "/import_key";
// const string wallet_sign_trx = wallet_func_base + "/sign_transaction";

func (api *API) GetAccount(name AccountName) (out *AccountResp, err error) {
	err = api.call("chain", "get_account", M{"account_name": name}, &out)
	return
}

func (api *API) GetCode(account AccountName) (out *Code, err error) {
	err = api.call("chain", "get_code", M{"account_name": account}, &out)
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
		hex.EncodeToString(api.ChainID), // eventually, we should receive the `chainID` from somewhere instead.
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

	resp, err := api.GetRequiredKeys(tx)
	if err != nil {
		return nil, fmt.Errorf("get_required_keys: %s", err)
	}

	stx := NewSignedTransaction(tx)

	stx.estimateResources(*opts, api.DefaultMaxCPUUsageMS, api.DefaultMaxNetUsageWords)

	signedTx, err := api.Signer.Sign(stx, api.ChainID, resp.RequiredKeys...)
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
