package eos

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// See: libraries/chain/include/eosio/chain/contracts/types.hpp:203
// See: build/contracts/eosio.system/eosio.system.abi

// Top-level actions

// SetCode represents the hard-coded `setcode` action.
type SetCode struct {
	Account   AccountName `json:"account"`
	VMType    byte        `json:"vmtype"`
	VMVersion byte        `json:"vmversion"`
	Code      HexBytes    `json:"bytes"`
}

// SetABI represents the hard-coded `setabi` action.
type SetABI struct {
	Account AccountName `json:"account"`
	ABI     ABI         `json:"abi"`
}

// belongs to `system`  structs
type EOSIOParameters struct {
	BasePerTransactionNetUsage     uint32 `json:"base_per_transaction_net_usage" yaml:"base_per_transaction_net_usage"`
	BasePerTransactionCPUUsage     uint32 `json:"base_per_transaction_cpu_usage" yaml:"base_per_transaction_cpu_usage"`
	BasePerActionCPUUsage          uint32 `json:"base_per_action_cpu_usage" yaml:"base_per_action_cpu_usage"`
	BaseSetcodeCPUUsage            uint32 `json:"base_setcode_cpu_usage" yaml:"base_setcode_cpu_usage"`
	PerSignatureCPUUsage           uint32 `json:"per_signature_cpu_usage" yaml:"per_signature_cpu_usage"`
	PerLockNetUsage                uint32 `json:"per_lock_net_usage" yaml:"per_lock_net_usage"`
	ContextFreeDiscountCPUUsageNum uint64 `json:"context_free_discount_cpu_usage_num" yaml:"context_free_discount_cpu_usage_num"`
	ContextFreeDiscountCPUUsageDen uint64 `json:"context_free_discount_cpu_usage_den" yaml:"context_free_discount_cpu_usage_den"`
	MaxTransactionCPUUsage         uint32 `json:"max_transaction_cpu_usage" yaml:"max_transaction_cpu_usage"`
	MaxTransactionNetUsage         uint32 `json:"max_transaction_net_usage" yaml:"max_transaction_net_usage"`

	MaxBlockCPUUsage       uint64 `json:"max_block_cpu_usage" yaml:"max_block_cpu_usage"`
	TargetBlockCPUUsagePct uint32 `json:"target_block_cpu_usage_pct" yaml:"target_block_cpu_usage_pct"` //< the target percent (1% == 100, 100%= 10,000) of maximum cpu usage; exceeding this triggers congestion handling
	MaxBblockNetUsage      uint64 `json:"max_block_net_usage" yaml:"max_block_net_usage"`               //< the maxiumum net usage in instructions for a block
	TargetBlockNetUsagePct uint32 `json:"target_block_net_usage_pct" yaml:"target_block_net_usage_pct"` //< the target percent (1% == 100, 100%= 10,000) of maximum net usage; exceeding this triggers congestion handling

	MaxTransactionLifetime       uint32 `json:"max_transaction_lifetime" yaml:"max_transaction_lifetime"`
	MaxTransactionExecTime       uint32 `json:"max_transaction_exec_time" yaml:"max_transaction_exec_time"`
	MaxAuthorityDepth            uint16 `json:"max_authority_depth" yaml:"max_authority_depth"`
	MaxInlineDepth               uint16 `json:"max_inline_depth" yaml:"max_inline_depth"`
	MaxInlineActionSize          uint32 `json:"max_inline_action_size" yaml:"max_inline_action_size"`
	MaxGeneratedTransactionCount uint32 `json:"max_generated_transaction_count" yaml:"max_generated_transaction_count"`

	// FIXME: does not appear in the `abi` for `eosio.system`.
	// MaxStorageSize uint64 `json:"max_storage_size" yaml:"max_storage_size"`
	PercentOfMaxInflationRate uint32 `json:"percent_of_max_inflation_rate" yaml:"percent_of_max_inflation_rate"`
	StorageReserveRatio       uint32 `json:"storage_reserve_ratio" yaml:"storage_reserve_ratio"`
}

// Sync with: /home/abourget/build/eos/patch1.patch

// belongs to the `system` structs
type EOSIOGlobalState struct {
	EOSIOParameters
	TotalStorageBytesReserved uint64 `json:"total_storage_bytes_reserved"`
	TotalStorageStake         uint64 `json:"total_storage_stake"`
	PaymentPerBlock           uint64 `json:"payment_per_block"`
}

// belongs to `system` structs
type DelegatedBandwidth struct {
	From         AccountName `json:"from"`
	To           AccountName `json:"to"`
	NetWeight    Asset       `json:"net_weight"`
	CPUWeight    Asset       `json:"cpu_weight"`
	StorageStake Asset       `json:"storage_stake"`
	StorageBytes uint64      `json:"storage_bytes"`
}

// belongs to `system` structs
type TotalResources struct {
}

// NewAccount represents the hard-coded `newaccount` action.
type NewAccount struct {
	Creator  AccountName `json:"creator"`
	Name     AccountName `json:"name"`
	Owner    Authority   `json:"owner"`
	Active   Authority   `json:"active"`
	Recovery Authority   `json:"recovery"`
}

// Action
type Action struct {
	Account       AccountName       `json:"account"`
	Name          ActionName        `json:"name"`
	Authorization []PermissionLevel `json:"authorization,omitempty"`
	Data          ActionData        `json:"data,omitempty"` // as HEX when we receive it.. FIXME: decode from hex directly.. and encode back plz!
}

type action struct {
	Account       AccountName       `json:"account"`
	Name          ActionName        `json:"name"`
	Authorization []PermissionLevel `json:"authorization,omitempty"`
}

type ActionData interface{}

func (a *Action) UnmarshalBinaryRead(r io.Reader) error {
	// Ok, we need to find a way to unmarshal those transactions.. to
	// do verification and/or introspection of blockchain data.
	//
	// There are two ways we can completely decode an incoming Action,
	// through a local map of structs (sort of a hard-coded ABI), or
	// through the ABI definitions and building an agnostic
	// map[string]interface{}.
	fmt.Println("MAMA")
	length, err := binary.ReadUvarint(&ByteReader{r})
	if err != nil {
		return err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}

	actionMap := registeredActions[a.Account]
	if actionMap == nil {
		return nil
	}

	objMap := actionMap[a.Name]
	if objMap == nil {
		return nil
	}

	obj := reflect.New(reflect.TypeOf(objMap))

	err = UnmarshalBinary(data, &obj)
	if err != nil {
		return err
	}

	a.Data = obj.Elem().Interface()

	return nil
}

var registeredActions = map[AccountName]map[ActionName]reflect.Type{}

// Registers Action objects..
func RegisterAction(accountName AccountName, actionName ActionName, obj interface{}) {
	// TODO: lock or som'th.. unless we never call after boot time..
	if registeredActions[accountName] == nil {
		registeredActions[accountName] = make(map[ActionName]reflect.Type)
	}
	registeredActions[accountName][actionName] = reflect.ValueOf(obj).Type()
}

func (a *Action) UnmarshalJSON(v []byte) (err error) {
	// load Account, Name, Authorization, Data
	// and then unpack other fields in a struct based on `Name` and `AccountName`..
	var newAct *action
	if err = json.Unmarshal(v, &newAct); err != nil {
		return
	}

	a.Account = newAct.Account
	a.Name = newAct.Name
	a.Authorization = newAct.Authorization

	actionMap := registeredActions[a.Account]
	if actionMap == nil {
		return nil
	}

	objMap := actionMap[a.Name]
	if objMap == nil {
		return nil
	}

	obj := reflect.New(reflect.TypeOf(objMap))
	err = json.Unmarshal(v, &obj)
	if err != nil {
		return err
	}

	a.Data = obj.Elem().Interface()

	return nil
}

func (a Action) MarshalBinary() ([]byte, error) {
	//fmt.Println("ENTERING MARSHALBINARY FOR ACTION!", a)
	// marshal binary a short action, then data
	common, err := MarshalBinary(&action{
		Account:       a.Account,
		Name:          a.Name,
		Authorization: a.Authorization,
	})
	if err != nil {
		return nil, err
	}

	var data []byte
	if a.Data != nil {
		data, err = MarshalBinary(a.Data)
		if err != nil {
			return nil, err
		}
	}

	varint := make([]byte, 4, 4)
	varintLen := binary.PutUvarint(varint, uint64(len(data)))
	// fmt.Println("****************************")
	// fmt.Println("ok mama", len(data), varint, varintLen)
	common = append(common, varint[:varintLen]...)
	// fmt.Println("ok mama", hex.EncodeToString(common))
	common = append(common, data...)
	// fmt.Println("ok mama", hex.EncodeToString(common))

	return common, nil
}

func (a *Action) MarshalJSON() ([]byte, error) {
	// Start with the base-line Action fields.

	cnt, err := json.Marshal(&action{
		Account:       a.Account,
		Name:          a.Name,
		Authorization: a.Authorization,
	})
	if err != nil {
		return nil, err
	}

	var keys1 map[string]interface{}
	err = json.Unmarshal(cnt, &keys1)
	if err != nil {
		return nil, err
	}

	if a.Data != nil {
		data, err := MarshalBinary(a.Data)
		if err != nil {
			return nil, err
		}

		keys1["data"] = hex.EncodeToString(data)
	}

	return json.Marshal(keys1)
}
