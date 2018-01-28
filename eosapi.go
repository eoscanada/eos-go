package eosapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type EOSAPI struct {
	HttpClient *http.Client
	BaseURL    string
}

func New(baseURL string) *EOSAPI {
	return &EOSAPI{
		HttpClient: http.DefaultClient,
		BaseURL:    baseURL,
	}
}

// Chain APIs
// Wallet APIs

func (api *EOSAPI) GetAccount(name AccountName) (out *AccountResp, err error) {
	err = api.call("POST", "chain", "get_account", M{"account_name": name}, &out)
	return
}

func (api *EOSAPI) GetCode(name AccountName) (out *Contract, err error) {
	err = api.call("POST", "chain", "get_code", M{"account_name": name}, &out)
	return
}

func (api *EOSAPI) call(method string, baseAPI string, endpoint string, body interface{}, out interface{}) error {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/v1/%s/%s", api.BaseURL, baseAPI, endpoint), enc(body))
	if err != nil {
		return err
	}

	resp, err := api.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var cnt bytes.Buffer
	_, err = io.Copy(&cnt, resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code = %d, body=%q", resp.StatusCode, cnt.String())
	}

	//fmt.Println("STRING", cnt.String())

	if err := json.Unmarshal(cnt.Bytes(), &out); err != nil {
		return err
	}

	return nil
}

type M map[string]interface{}

func enc(v interface{}) io.Reader {
	if v == nil {
		return nil
	}

	cnt, _ := json.Marshal(v)

	return bytes.NewReader(cnt)
}
