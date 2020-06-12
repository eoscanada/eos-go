package boot

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type Snapshot []SnapshotLine

type SnapshotLine struct {
	EthereumAddress string
	EOSPublicKey    ecc.PublicKey
	Balance         eos.Asset
	AccountName     string
}

func NewSnapshot(content []byte) (out Snapshot, err error) {
	reader := csv.NewReader(bytes.NewBuffer(content))
	allRecords, err := reader.ReadAll()
	if err != nil {
		return
	}

	for _, el := range allRecords {
		if len(el) != 4 {
			return nil, fmt.Errorf("should have 4 elements per line")
		}

		newAsset, err := eos.NewEOSAssetFromString(el[3])
		if err != nil {
			return out, err
		}

		pubKey, err := ecc.NewPublicKey(el[2])
		if err != nil {
			return out, err
		}

		out = append(out, SnapshotLine{el[0], pubKey, newAsset, el[1]})
	}

	return
}

type UnregdSnapshot []UnregdSnapshotLine

type UnregdSnapshotLine struct {
	EthereumAddress string
	AccountName     string
	Balance         eos.Asset
}

func NewUnregdSnapshot(content []byte) (out UnregdSnapshot, err error) {
	reader := csv.NewReader(bytes.NewBuffer(content))
	allRecords, err := reader.ReadAll()
	if err != nil {
		return
	}

	for _, el := range allRecords {
		if len(el) != 3 {
			return nil, fmt.Errorf("should have 2 elements per line")
		}

		newAsset, err := eos.NewEOSAssetFromString(el[2])
		if err != nil {
			return out, err
		}

		out = append(out, UnregdSnapshotLine{el[0], el[1], newAsset})
	}

	return
}
