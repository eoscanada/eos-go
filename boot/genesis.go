package boot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// TODO: update with latest GenesisJSON with the basic parameters...
type GenesisJSON struct {
	InitialTimestamp string `json:"initial_timestamp"`
	InitialKey       string `json:"initial_key"`
}

func generateGenesisJSON(pubKey string) string {
	// known not to fail
	cnt, _ := json.Marshal(&GenesisJSON{
		InitialTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05"),
		InitialKey:       pubKey,
	})
	return string(cnt)
}

func loadGenesisFromFile(pubkey string) (string, error) {
	cnt, err := ioutil.ReadFile("genesis.json")
	if err != nil {
		return "", err
	}

	var gendata *GenesisJSON
	err = json.Unmarshal(cnt, &gendata)
	if err != nil {
		return "", err
	}

	if pubkey != gendata.InitialKey {
		return "", fmt.Errorf("attempting to reuse genesis.json: genesis.key doesn't match genesis.json")
	}

	out, _ := json.Marshal(gendata)

	return string(out), nil
}
