package eos

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"

	"github.com/eoscanada/eos-go/ecc"
)

type API struct {
	HttpClient *http.Client
	BaseURL    string
	Signer     Signer
	Debug      bool
	Compress   CompressionType
	// Header is one or more headers to be added to all outgoing calls
	Header                  http.Header
	DefaultMaxCPUUsageMS    uint8
	DefaultMaxNetUsageWords uint32 // in 8-bytes words

	lastGetInfo      *InfoResp
	lastGetInfoStamp time.Time
	lastGetInfoLock  sync.Mutex

	customGetRequiredKeys     func(ctx context.Context, tx *Transaction) ([]ecc.PublicKey, error)
	enablePartialRequiredKeys bool
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
		BaseURL:  strings.TrimRight(baseURL, "/"),
		Compress: CompressionZlib,
		Header:   make(http.Header),
	}

	return api
}

// FixKeepAlives tests the remote server for keepalive support (the
// main `nodeos` software doesn't in the version from March 22nd
// 2018).  Some endpoints front their node with a keep-alive
// supporting web server.  Adjust the `KeepAlive` support of the
// client accordingly.
func (api *API) FixKeepAlives(ctx context.Context) bool {
	// Yeah, to provoke a keep alive, you need to query twice.
	for i := 0; i < 5; i++ {
		_, err := api.GetInfo(ctx)
		if api.Debug {
			log.Println("err", err)
		}
		if err == io.EOF {
			if tr, ok := api.HttpClient.Transport.(*http.Transport); ok {
				tr.DisableKeepAlives = true
				return true
			}
		}
		_, err = api.GetNetConnections(ctx)
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

func (api *API) SetCustomGetRequiredKeys(f func(ctx context.Context, tx *Transaction) ([]ecc.PublicKey, error)) {
	api.customGetRequiredKeys = f
}

func (api *API) UsePartialRequiredKeys() {
	api.enablePartialRequiredKeys = true
}

func (api *API) getPartialRequiredKeys(ctx context.Context, tx *Transaction) ([]ecc.PublicKey, error) {
	// loop to get all the authorizers, and dedupe
	var authorizers []PermissionLevel
	for _, act := range tx.Actions {
		for _, pl := range act.Authorization {
			authorizers = append(authorizers, pl)
		}
	}

	ourKeys, err := api.Signer.AvailableKeys(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := api.GetAccountsByAuthorizers(ctx, authorizers, ourKeys)
	if err != nil {
		return nil, err
	}

	var out []ecc.PublicKey
	seen := make(map[string]bool)
	for _, acct := range resp.Accounts {
		if acct.AuthorizingKey == nil {
			continue
		}

		stringKey := acct.AuthorizingKey.String()
		if seen[stringKey] {
			continue
		}

		out = append(out, *acct.AuthorizingKey)
		seen[stringKey] = true
	}

	return out, nil
}

func (api *API) SetSigner(s Signer) {
	api.Signer = s
}

// ProducerPause will pause block production on a nodeos with
// `producer_api` plugin loaded.
func (api *API) ProducerPause(ctx context.Context) error {
	return api.call(ctx, "producer", "pause", nil, nil)
}

// CreateSnapshot will write a snapshot file on a nodeos with
// `producer_api` plugin loaded.
func (api *API) CreateSnapshot(ctx context.Context) (out *CreateSnapshotResp, err error) {
	err = api.call(ctx, "producer", "create_snapshot", nil, &out)
	return
}

// GetIntegrityHash will produce a hash corresponding to current
// state. Requires `producer_api` and useful when loading
// from a snapshot
func (api *API) GetIntegrityHash(ctx context.Context) (out *GetIntegrityHashResp, err error) {
	err = api.call(ctx, "producer", "get_integrity_hash", nil, &out)
	return
}

// ProducerResume will resume block production on a nodeos with
// `producer_api` plugin loaded. Obviously, this needs to be a
// producing node on the producers schedule for it to do anything.
func (api *API) ProducerResume(ctx context.Context) error {
	return api.call(ctx, "producer", "resume", nil, nil)
}

// IsProducerPaused queries the blockchain for the pause statement of
// block production.
func (api *API) IsProducerPaused(ctx context.Context) (out bool, err error) {
	err = api.call(ctx, "producer", "paused", nil, &out)
	return
}

func (api *API) GetProducerProtocolFeatures(ctx context.Context) (out []ProtocolFeature, err error) {
	err = api.call(ctx, "producer", "get_supported_protocol_features", nil, &out)
	return
}

func (api *API) ScheduleProducerProtocolFeatureActivations(ctx context.Context, protocolFeaturesToActivate []Checksum256) error {
	return api.call(ctx, "producer", "schedule_protocol_feature_activations", M{"protocol_features_to_activate": protocolFeaturesToActivate}, nil)
}

func (api *API) GetAccount(ctx context.Context, name AccountName) (out *AccountResp, err error) {
	err = api.call(ctx, "chain", "get_account", M{"account_name": name}, &out)
	return
}

func (api *API) GetRawCodeAndABI(ctx context.Context, account AccountName) (out *GetRawCodeAndABIResp, err error) {
	err = api.call(ctx, "chain", "get_raw_code_and_abi", M{"account_name": account}, &out)
	return
}

func (api *API) GetCode(ctx context.Context, account AccountName) (out *GetCodeResp, err error) {
	err = api.call(ctx, "chain", "get_code", M{"account_name": account, "code_as_wasm": true}, &out)
	return
}

func (api *API) GetCodeHash(ctx context.Context, account AccountName) (out Checksum256, err error) {
	resp := GetCodeHashResp{}
	if err = api.call(ctx, "chain", "get_code_hash", M{"account_name": account}, &resp); err != nil {
		return
	}

	buffer, err := hex.DecodeString(resp.CodeHash)
	return Checksum256(buffer), err
}

func (api *API) GetABI(ctx context.Context, account AccountName) (out *GetABIResp, err error) {
	err = api.call(ctx, "chain", "get_abi", M{"account_name": account}, &out)
	return
}

func (api *API) ABIJSONToBin(ctx context.Context, code AccountName, action Name, payload M) (out HexBytes, err error) {
	resp := ABIJSONToBinResp{}
	err = api.call(ctx, "chain", "abi_json_to_bin", M{"code": code, "action": action, "args": payload}, &resp)
	if err != nil {
		return
	}

	buffer, err := hex.DecodeString(resp.Binargs)
	return HexBytes(buffer), err
}

func (api *API) ABIBinToJSON(ctx context.Context, code AccountName, action Name, payload HexBytes) (out M, err error) {
	resp := ABIBinToJSONResp{}
	err = api.call(ctx, "chain", "abi_bin_to_json", M{"code": code, "action": action, "binargs": payload}, &resp)
	if err != nil {
		return
	}

	return resp.Args, nil
}

func (api *API) WalletCreate(ctx context.Context, walletName string) (err error) {
	return api.call(ctx, "wallet", "create", walletName, nil)
}

func (api *API) WalletOpen(ctx context.Context, walletName string) (err error) {
	return api.call(ctx, "wallet", "open", walletName, nil)
}

func (api *API) WalletLock(ctx context.Context, walletName string) (err error) {
	return api.call(ctx, "wallet", "lock", walletName, nil)
}

func (api *API) WalletLockAll(ctx context.Context) (err error) {
	return api.call(ctx, "wallet", "lock_all", nil, nil)
}

func (api *API) WalletUnlock(ctx context.Context, walletName, password string) (err error) {
	return api.call(ctx, "wallet", "unlock", []string{walletName, password}, nil)
}

// WalletImportKey loads a new WIF-encoded key into the wallet.
func (api *API) WalletImportKey(ctx context.Context, walletName, wifPrivKey string) (err error) {
	return api.call(ctx, "wallet", "import_key", []string{walletName, wifPrivKey}, nil)
}

func (api *API) WalletPublicKeys(ctx context.Context) (out []ecc.PublicKey, err error) {
	var textKeys []string
	err = api.call(ctx, "wallet", "get_public_keys", nil, &textKeys)
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

func (api *API) ListWallets(ctx context.Context, walletName ...string) (out []string, err error) {
	err = api.call(ctx, "wallet", "list_wallets", walletName, &out)
	if err != nil {
		return nil, err
	}

	return
}

func (api *API) ListKeys(ctx context.Context, walletNames ...string) (out []*ecc.PrivateKey, err error) {
	var textKeys []string
	err = api.call(ctx, "wallet", "list_keys", walletNames, &textKeys)
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

func (api *API) GetPublicKeys(ctx context.Context) (out []*ecc.PublicKey, err error) {
	var textKeys []string
	err = api.call(ctx, "wallet", "get_public_keys", nil, &textKeys)
	if err != nil {
		return nil, err
	}

	for _, k := range textKeys {
		newKey, err := ecc.NewPublicKey(k)
		if err != nil {
			return nil, err
		}

		out = append(out, &newKey)
	}
	return
}

func (api *API) WalletSetTimeout(ctx context.Context, timeout int32) (err error) {
	return api.call(ctx, "wallet", "set_timeout", timeout, nil)
}

func (api *API) WalletSignTransaction(ctx context.Context, tx *SignedTransaction, chainID []byte, pubKeys ...ecc.PublicKey) (out *WalletSignTransactionResp, err error) {
	var textKeys []string
	for _, key := range pubKeys {
		textKeys = append(textKeys, key.String())
	}

	err = api.call(ctx, "wallet", "sign_transaction", []interface{}{
		tx,
		textKeys,
		hex.EncodeToString(chainID),
	}, &out)
	return
}

// SignPushActions will create a transaction, fill it with default
// values, sign it and submit it to the chain.  It is the highest
// level function on top of the `/v1/chain/push_transaction` endpoint.
func (api *API) SignPushActions(ctx context.Context, a ...*Action) (out *PushTransactionFullResp, err error) {
	return api.SignPushActionsWithOpts(ctx, a, nil)
}

func (api *API) SignPushActionsWithOpts(ctx context.Context, actions []*Action, opts *TxOptions) (out *PushTransactionFullResp, err error) {
	if opts == nil {
		opts = &TxOptions{}
	}

	if err := opts.FillFromChain(ctx, api); err != nil {
		return nil, err
	}

	tx := NewTransaction(actions, opts)

	return api.SignPushTransaction(ctx, tx, opts.ChainID, opts.Compress)
}

// SignPushTransaction will sign a transaction and submit it to the
// chain.
func (api *API) SignPushTransaction(ctx context.Context, tx *Transaction, chainID Checksum256, compression CompressionType) (out *PushTransactionFullResp, err error) {
	_, packed, err := api.SignTransaction(ctx, tx, chainID, compression)
	if err != nil {
		return nil, err
	}

	return api.PushTransaction(ctx, packed)
}

// SignTransaction will sign and pack a transaction, but not submit to
// the chain.  It lives on the `api` object because it might query the
// blockchain to learn which keys are required to sign this particular
// transaction.
//
// You can override the way we request keys (which defaults to
// `api.GetRequiredKeys()`) with SetCustomGetRequiredKeys().
//
// To sign a transaction, you need a Signer defined on the `API`
// object. See SetSigner.
func (api *API) SignTransaction(ctx context.Context, tx *Transaction, chainID Checksum256, compression CompressionType) (*SignedTransaction, *PackedTransaction, error) {
	if api.Signer == nil {
		return nil, nil, fmt.Errorf("no Signer configured")
	}

	stx := NewSignedTransaction(tx)

	var err error
	var requiredKeys []ecc.PublicKey
	if api.customGetRequiredKeys != nil {
		requiredKeys, err = api.customGetRequiredKeys(ctx, tx)
		if err != nil {
			return nil, nil, fmt.Errorf("custom_get_required_keys: %s", err)
		}
	} else {
		if api.enablePartialRequiredKeys {
			requiredKeys, err = api.getPartialRequiredKeys(ctx, tx)
			if err != nil {
				return nil, nil, fmt.Errorf("get_accounts_by_authorizers: %s", err)
			}
		} else {
			resp, err := api.GetRequiredKeys(ctx, tx)
			if err != nil {
				return nil, nil, fmt.Errorf("get_required_keys: %s", err)
			}
			requiredKeys = resp.RequiredKeys
		}
	}

	signedTx, err := api.Signer.Sign(ctx, stx, chainID, requiredKeys...)
	if err != nil {
		return nil, nil, fmt.Errorf("signing through wallet: %s", err)
	}

	packed, err := signedTx.Pack(compression)
	if err != nil {
		return nil, nil, err
	}

	return signedTx, packed, nil
}

// PushTransaction submits a properly filled (tapos), packed and
// signed transaction to the blockchain.
func (api *API) PushTransaction(ctx context.Context, tx *PackedTransaction) (out *PushTransactionFullResp, err error) {
	err = api.call(ctx, "chain", "push_transaction", tx, &out)
	return
}

func (api *API) SendTransaction(ctx context.Context, tx *PackedTransaction) (out *PushTransactionFullResp, err error) {
	err = api.call(ctx, "chain", "send_transaction", tx, &out)
	return
}

func (api *API) PushTransactionRaw(ctx context.Context, tx *PackedTransaction) (out json.RawMessage, err error) {
	err = api.call(ctx, "chain", "push_transaction", tx, &out)
	return
}
func (api *API) SendTransactionRaw(ctx context.Context, tx *PackedTransaction) (out json.RawMessage, err error) {
	err = api.call(ctx, "chain", "send_transaction", tx, &out)
	return
}

func (api *API) GetInfo(ctx context.Context) (out *InfoResp, err error) {
	err = api.call(ctx, "chain", "get_info", nil, &out)
	return
}

func (api *API) cachedGetInfo(ctx context.Context) (*InfoResp, error) {
	api.lastGetInfoLock.Lock()
	defer api.lastGetInfoLock.Unlock()

	var info *InfoResp
	var err error

	if !api.lastGetInfoStamp.IsZero() && time.Now().Add(-1*time.Second).Before(api.lastGetInfoStamp) {
		info = api.lastGetInfo
	} else {
		info, err = api.GetInfo(ctx)
		if err != nil {
			return nil, err
		}
		api.lastGetInfoStamp = time.Now()
		api.lastGetInfo = info
	}
	if err != nil {
		return nil, err
	}

	return api.lastGetInfo, nil
}

func (api *API) GetNetConnections(ctx context.Context) (out []*NetConnectionsResp, err error) {
	err = api.call(ctx, "net", "connections", nil, &out)
	return
}

func (api *API) NetConnect(ctx context.Context, host string) (out NetConnectResp, err error) {
	err = api.call(ctx, "net", "connect", host, &out)
	return
}

func (api *API) NetDisconnect(ctx context.Context, host string) (out NetDisconnectResp, err error) {
	err = api.call(ctx, "net", "disconnect", host, &out)
	return
}

func (api *API) GetNetStatus(ctx context.Context, host string) (out *NetStatusResp, err error) {
	err = api.call(ctx, "net", "status", M{"host": host}, &out)
	return
}

func (api *API) GetBlockByID(ctx context.Context, id string) (out *BlockResp, err error) {
	err = api.call(ctx, "chain", "get_block", M{"block_num_or_id": id}, &out)
	return
}

// GetScheduledTransactionsWithBounds returns scheduled transactions within specified bounds
func (api *API) GetScheduledTransactionsWithBounds(ctx context.Context, lower_bound string, limit uint32) (out *ScheduledTransactionsResp, err error) {
	err = api.call(ctx, "chain", "get_scheduled_transactions", M{"json": true, "lower_bound": lower_bound, "limit": limit}, &out)
	return
}

// GetScheduledTransactions returns the Top 100 scheduled transactions
func (api *API) GetScheduledTransactions(ctx context.Context) (out *ScheduledTransactionsResp, err error) {
	return api.GetScheduledTransactionsWithBounds(ctx, "", 100)
}

func (api *API) GetProducers(ctx context.Context) (out *ProducersResp, err error) {
	/*
		+FC_REFLECT( eosio::chain_apis::read_only::get_producers_params, (json)(lower_bound)(limit) )
		+FC_REFLECT( eosio::chain_apis::read_only::get_producers_result, (rows)(total_producer_vote_weight)(more) ); */
	err = api.call(ctx, "chain", "get_producers", nil, &out)
	return
}

func (api *API) GetBlockByNum(ctx context.Context, num uint32) (out *BlockResp, err error) {
	err = api.call(ctx, "chain", "get_block", M{"block_num_or_id": fmt.Sprintf("%d", num)}, &out)
	//err = api.call("chain", "get_block", M{"block_num_or_id": num}, &out)
	return
}

func (api *API) GetBlockByNumOrID(ctx context.Context, query string) (out *SignedBlock, err error) {
	err = api.call(ctx, "chain", "get_block", M{"block_num_or_id": query}, &out)
	return
}

func (api *API) GetBlockByNumOrIDRaw(ctx context.Context, query string) (out interface{}, err error) {
	err = api.call(ctx, "chain", "get_block", M{"block_num_or_id": query}, &out)
	return
}

func (api *API) GetDBSize(ctx context.Context) (out *DBSizeResp, err error) {
	err = api.call(ctx, "db_size", "get", nil, &out)
	return
}

func (api *API) GetTransaction(ctx context.Context, id string) (out *TransactionResp, err error) {
	err = api.call(ctx, "history", "get_transaction", M{"id": id}, &out)
	return
}

func (api *API) GetTransactionRaw(ctx context.Context, id string) (out json.RawMessage, err error) {
	err = api.call(ctx, "history", "get_transaction", M{"id": id}, &out)
	return
}

func (api *API) GetActions(ctx context.Context, params GetActionsRequest) (out *ActionsResp, err error) {
	err = api.call(ctx, "history", "get_actions", params, &out)
	return
}

func (api *API) GetKeyAccounts(ctx context.Context, publicKey string) (out *KeyAccountsResp, err error) {
	err = api.call(ctx, "history", "get_key_accounts", M{"public_key": publicKey}, &out)
	return
}

func (api *API) GetControlledAccounts(ctx context.Context, controllingAccount string) (out *ControlledAccountsResp, err error) {
	err = api.call(ctx, "history", "get_controlled_accounts", M{"controlling_account": controllingAccount}, &out)
	return
}

func (api *API) GetTransactions(ctx context.Context, name AccountName) (out *TransactionsResp, err error) {
	err = api.call(ctx, "account_history", "get_transactions", M{"account_name": name}, &out)
	return
}

func (api *API) GetTableByScope(ctx context.Context, params GetTableByScopeRequest) (out *GetTableByScopeResp, err error) {
	err = api.call(ctx, "chain", "get_table_by_scope", params, &out)
	return
}

func (api *API) GetTableRows(ctx context.Context, params GetTableRowsRequest) (out *GetTableRowsResp, err error) {
	err = api.call(ctx, "chain", "get_table_rows", params, &out)
	return
}

func (api *API) GetRawABI(ctx context.Context, params GetRawABIRequest) (out *GetRawABIResp, err error) {
	err = api.call(ctx, "chain", "get_raw_abi", params, &out)
	return
}

func (api *API) GetRequiredKeys(ctx context.Context, tx *Transaction) (out *GetRequiredKeysResp, err error) {
	keys, err := api.Signer.AvailableKeys(ctx)
	if err != nil {
		return nil, err
	}

	err = api.call(ctx, "chain", "get_required_keys", M{"transaction": tx, "available_keys": keys}, &out)
	return
}

func (api *API) GetAccountsByAuthorizers(ctx context.Context, authorizations []PermissionLevel, keys []ecc.PublicKey) (out *GetAccountsByAuthorizersResp, err error) {
	err = api.call(ctx, "chain", "get_accounts_by_authorizers", M{"accounts": authorizations, "keys": keys}, &out)
	return
}

func (api *API) GetCurrencyBalance(ctx context.Context, account AccountName, symbol string, code AccountName) (out []Asset, err error) {
	params := M{"account": account, "code": code}
	if symbol != "" {
		params["symbol"] = symbol
	}
	err = api.call(ctx, "chain", "get_currency_balance", params, &out)
	return
}

func (api *API) GetCurrencyStats(ctx context.Context, code AccountName, symbol string) (out *GetCurrencyStatsResp, err error) {
	params := M{"code": code, "symbol": symbol}

	outWrapper := make(map[string]*GetCurrencyStatsResp)
	err = api.call(ctx, "chain", "get_currency_stats", params, &outWrapper)
	out = outWrapper[symbol]

	return
}

// See more here: libraries/chain/contracts/abi_serializer.cpp:58...

func (api *API) call(ctx context.Context, baseAPI string, endpoint string, body interface{}, out interface{}) error {
	jsonBody, err := enc(body)
	if err != nil {
		return err
	}

	targetURL := fmt.Sprintf("%s/v1/%s/%s", api.BaseURL, baseAPI, endpoint)
	req, err := http.NewRequest("POST", targetURL, jsonBody)
	if err != nil {
		return fmt.Errorf("NewRequest: %s", err)
	}

	for k, v := range api.Header {
		if req.Header == nil {
			req.Header = http.Header{}
		}
		req.Header[k] = append(req.Header[k], v...)
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

	resp, err := api.HttpClient.Do(req.WithContext(ctx))
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
		var apiErr APIError
		if err := json.Unmarshal(cnt.Bytes(), &apiErr); err != nil {
			return ErrNotFound
		}
		return apiErr
	}

	if resp.StatusCode > 299 {
		var apiErr APIError
		if err := json.Unmarshal(cnt.Bytes(), &apiErr); err != nil {
			return fmt.Errorf("%s: status code=%d, body=%s", req.URL.String(), resp.StatusCode, cnt.String())
		}

		// Handle cases where some API calls (/v1/chain/get_account for example) returns a 500
		// error when retrieving data that does not exist.
		if apiErr.IsUnknownKeyError() {
			return ErrNotFound
		}

		return apiErr
	}

	if api.Debug {
		fmt.Println("RESPONSE:")
		responseDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("-------------------------------")
		fmt.Println(cnt.String())
		fmt.Println("-------------------------------")
		fmt.Printf("%q\n", responseDump)
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

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
