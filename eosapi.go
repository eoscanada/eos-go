package eosapi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/eosioca/eosapi/ecc"
)

type EOSAPI struct {
	HttpClient *http.Client
	BaseURL    *url.URL
	ChainID    []byte
	Signer     Signer
}

func New(baseURL *url.URL, chainID []byte) *EOSAPI {
	if len(chainID) != 32 {
		panic("chainID must be 32 bytes")
	}

	api := &EOSAPI{
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
		BaseURL: baseURL,
		ChainID: chainID,
	}

	return api
}

// FixKeepAlives tests the remote server for keepalive support (the
// main `nodeos` software doesn't in the version from March 22nd
// 2018).  Some endpoints front their node with a keep-alive
// supporting web server.  Adjust the `KeepAlive` support of the
// client accordingly.
func (api *EOSAPI) FixKeepAlives() bool {
	// Yeah, to provoke a keep alive, you need to query twice.
	for i := 0; i < 5; i++ {
		_, err := api.GetInfo()
		log.Println("err", err)
		if err == io.EOF {
			if tr, ok := api.HttpClient.Transport.(*http.Transport); ok {
				tr.DisableKeepAlives = true
				return true
			}
		}
		_, err = api.GetNetConnections()
		log.Println("err", err)
		if err == io.EOF {
			if tr, ok := api.HttpClient.Transport.(*http.Transport); ok {
				tr.DisableKeepAlives = true
				return true
			}
		}
	}
	return false
}

func (api *EOSAPI) EnableKeepAlives() bool {
	if tr, ok := api.HttpClient.Transport.(*http.Transport); ok {
		tr.DisableKeepAlives = false
		return true
	}
	return false
}

func (api *EOSAPI) SetSigner(s Signer) {
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

func (api *EOSAPI) GetAccount(name AccountName) (out *AccountResp, err error) {
	err = api.call("chain", "get_account", M{"account_name": name}, &out)
	return
}

func (api *EOSAPI) GetCode(account AccountName) (out *Code, err error) {
	err = api.call("chain", "get_code", M{"account_name": account}, &out)
	return
}

func (api *EOSAPI) WalletPublicKeys() (out []ecc.PublicKey, err error) {
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

func (api *EOSAPI) WalletSignTransaction(tx *SignedTransaction, pubKeys ...ecc.PublicKey) (out *WalletSignTransactionResp, err error) {
	var textKeys []string
	for _, key := range pubKeys {
		textKeys = append(textKeys, key.String())
	}

	err = api.call("wallet", "sign_transaction", []interface{}{
		tx,
		textKeys,
		hex.EncodeToString(api.ChainID),
	}, &out)
	return
}

func (api *EOSAPI) PushSignedTransaction(tx *SignedTransaction) (out *PushTransactionFullResp, err error) {

	//fmt.Println("PUSHING signed transaction", tx.Transaction)
	data, err := MarshalBinary(tx.Transaction)
	if err != nil {
		return nil, err
	}

	//fmt.Println("hex data", hex.EncodeToString(data))

	err = api.call("chain", "push_transaction", M{"data": hex.EncodeToString(data), "signatures": tx.Signatures, "compression": "none"}, &out)
	return
}

func (api *EOSAPI) NewAccount(creator, newAccount AccountName, publicKey ecc.PublicKey) (out *PushTransactionFullResp, err error) {
	if api.Signer == nil {
		return nil, fmt.Errorf("no Signer configured")
	}

	a := &Action{
		Account: AccountName("eosio"),
		Name:    ActionName("newaccount"),
		Authorization: []PermissionLevel{
			{creator, PermissionName("active")},
		},
		Data: NewAccount{
			Creator: creator,
			Name:    newAccount,
			Owner: Authority{
				Threshold: 1,
				Keys: []KeyWeight{
					KeyWeight{
						PublicKey: publicKey,
						Weight:    1,
					},
				},
			},
			Active: Authority{
				Threshold: 1,
				Keys: []KeyWeight{
					KeyWeight{
						PublicKey: publicKey,
						Weight:    1,
					},
				},
			},
			Recovery: Authority{
				Threshold: 1,
				Accounts: []PermissionLevelWeight{
					PermissionLevelWeight{
						Permission: PermissionLevel{creator, PermissionName("active")},
						Weight:     1,
					},
				},
			},
		},
	}
	tx := &Transaction{
		Actions: []*Action{a},
	}

	chainID, err := tx.Fill(api)
	if err != nil {
		return nil, err
	}

	resp, err := api.GetRequiredKeys(tx, api.Signer)
	log.Println("GetRequiredKeys", resp, err)
	if err != nil {
		return nil, fmt.Errorf("GetRequiredKeys: %s", err)
	}

	signedTx, err := api.Signer.Sign(NewSignedTransaction(tx), chainID, resp.RequiredKeys...)
	if err != nil {
		return nil, fmt.Errorf("Sign: %s", err)
	}

	return api.PushSignedTransaction(signedTx)
}

func (api *EOSAPI) SetCode(account AccountName, wasmPath, abiPath string) (out *PushTransactionFullResp, err error) {
	codeContent, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return nil, err
	}

	tx := &Transaction{
		Actions: []*Action{
			{
				Account: AccountName("eosio"),
				Name:    ActionName("setcode"),
				Authorization: []PermissionLevel{
					{account, PermissionName("active")},
				},
				Data: SetCode{
					Account:   account,
					VMType:    0,
					VMVersion: 0,
					Code:      HexBytes(codeContent),
				},
			},
			{
				Account: AccountName("eosio"),
				Name:    ActionName("setabi"),
				Authorization: []PermissionLevel{
					{account, PermissionName("active")},
				},
				Data: SetCode{
					Account:   account,
					VMType:    0,
					VMVersion: 0,
					Code:      HexBytes(codeContent),
				},
			},
		},
	}
	tx = &Transaction{
		Actions: []*Action{
			{
				Account: AccountName("eosio"),
				Name:    ActionName("issue"),
				Authorization: []PermissionLevel{
					{AccountName("eosio"), PermissionName("active")},
				},
				Data: Issue{
					To:       AccountName("abourget"),
					Quantity: 123123,
				},
			},
		},
	}

	chainID, err := tx.Fill(api)
	if err != nil {
		return nil, err
	}

	resp, err := api.GetRequiredKeys(tx, api.Signer)
	log.Println("GetRequiredKeys", resp, err)
	if err != nil {
		return nil, fmt.Errorf("GetRequiredKeys: %s", err)
	}

	signed, err := api.Signer.Sign(NewSignedTransaction(tx), chainID, resp.RequiredKeys...)
	if err != nil {
		return nil, fmt.Errorf("Sign: %s", err)
	}

	return api.PushSignedTransaction(signed)
}

func (api *EOSAPI) GetInfo() (out *InfoResp, err error) {
	err = api.call("chain", "get_info", nil, &out)
	return
}

func (api *EOSAPI) GetNetConnections() (out []*NetConnectionsResp, err error) {
	err = api.call("net", "connections", nil, &out)
	return
}

func (api *EOSAPI) NetConnect(host string) (out NetConnectResp, err error) {
	err = api.call("net", "connect", host, &out)
	return
}

func (api *EOSAPI) NetDisconnect(host string) (out NetDisconnectResp, err error) {
	err = api.call("net", "disconnect", host, &out)
	return
}

func (api *EOSAPI) GetNetStatus(host string) (out *NetStatusResp, err error) {
	err = api.call("net", "status", M{"host": host}, &out)
	return
}

func (api *EOSAPI) GetBlockByID(id string) (out *BlockResp, err error) {
	err = api.call("chain", "get_block", M{"block_num_or_id": id}, &out)
	return
}

func (api *EOSAPI) GetBlockByNum(num uint64) (out *BlockResp, err error) {
	err = api.call("chain", "get_block", M{"block_num_or_id": fmt.Sprintf("%d", num)}, &out)
	return
}

func (api *EOSAPI) GetTableRows(params GetTableRowsRequest) (out *GetTableRowsResp, err error) {
	err = api.call("chain", "get_table_rows", params, &out)
	return
}

func (api *EOSAPI) GetRequiredKeys(tx *Transaction, signer Signer) (out *GetRequiredKeysResp, err error) {
	keys, err := signer.AvailableKeys()
	if err != nil {
		return nil, err
	}

	err = api.call("chain", "get_required_keys", M{"transaction": tx, "available_keys": keys}, &out)
	return
}

func (api *EOSAPI) GetCurrencyBalance(account AccountName, symbol string, code AccountName) (out []Asset, err error) {
	err = api.call("chain", "get_currency_balance", M{"account": account, "symbol": symbol, "code": code}, &out)
	return
}

// See more here: libraries/chain/contracts/abi_serializer.cpp:58...

func (api *EOSAPI) call(baseAPI string, endpoint string, body interface{}, out interface{}) error {
	jsonBody, err := enc(body)
	if err != nil {
		return err
	}

	baseURL := &url.URL{
		Scheme: api.BaseURL.Scheme,
		Host:   api.BaseURL.Host,
		Path:   fmt.Sprintf("/v1/%s/%s", baseAPI, endpoint),
	}
	req, err := http.NewRequest("POST", baseURL.String(), jsonBody)
	if err != nil {
		return fmt.Errorf("NewRequest: %s", err)
	}

	// Useful when debugging API calls
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	resp, err := api.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Do: %s", err)
	}
	defer resp.Body.Close()

	var cnt bytes.Buffer
	_, err = io.Copy(&cnt, resp.Body)
	if err != nil {
		return fmt.Errorf("Copy: %s", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("status code=%d, body=%s", resp.StatusCode, cnt.String())
	}

	fmt.Println("SERVER RESPONSE", cnt.String())

	if err := json.Unmarshal(cnt.Bytes(), &out); err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}

	return nil
}

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
