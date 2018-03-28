package eosapi

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"reflect"
)

// See: libraries/chain/include/eosio/chain/contracts/types.hpp:203
// See: build/contracts/eosio.system/eosio.system.abi

// Top-level actions

// Transfer represents the `eosio.system::transfer` action.
type Transfer struct {
	From     AccountName `json:"from"`
	To       AccountName `json:"to"`
	Quantity uint64      `json:"quantity"`
	Memo     string      `json:"memo"`
}

// Issue represents the `eosio.system::issue` action.
type Issue struct {
	To       AccountName `json:"to"`
	Quantity uint64      `json:"quantity" struc:"uint64,little"`
}

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

// DelegateBW represents the `eosio.system::delegatebw` action.
type DelegateBW struct {
	From         AccountName `json:"from"`
	Receiver     AccountName `json:"receiver"`
	StakeNet     Asset       `json:"stake_net"`
	StakeCPU     Asset       `json:"stake_cpu"`
	StakeStorage Asset       `json:"stake_storage"`
}

// UndelegateBW represents the `eosio.system::undelegatebw` action.
type UndelegateBW struct {
	From         AccountName `json:"from"`
	Receiver     AccountName `json:"receiver"`
	UnstakeNet   Asset       `json:"unstake_net"`
	UnstakeCPU   Asset       `json:"unstake_cpu"`
	UnstakeBytes uint64      `json:"unstake_bytes"`
}

// Refund represents the `eosio.system::refund` action
type Refund struct {
	Owner AccountName `json:"owner"`
}

// RegisterProducer represents the `eosio.system::regproducer` action
type RegisterProducer struct {
	Producer    AccountName     `json:"producer"`
	ProducerKey []byte          `json:"producer_key"`
	Prefs       EOSIOParameters `json:"eosio_parameters"`
}

// UnregisterProducer represents the `eosio.system::unregprod` action
type UnregisterProducer struct {
	Producer AccountName `json:"producer"`
}

// RegisterProxy represents the `eosio.system::regproxy` action
type RegisterProxy struct {
	Proxy AccountName `json:"proxy"`
}

// UnregisterProxy represents the `eosio.system::unregproxy` action
type UnregisterProxy struct {
	Proxy AccountName `json:"proxy"`
}

// VoteProducer represents the `eosio.system::voteproducer` action
type VoteProducer struct {
	Voter     AccountName   `json:"voter"`
	Proxy     AccountName   `json:"proxy"`
	Producers []AccountName `json:"producers"`
}

// ClaimRewards repreents the `eosio.system::claimrewards` action
type ClaimRewards struct {
	Owner AccountName `json:"owner"`
}

// Nonce represents the `eosio.system::nonce` action. It is used to
// add variability in a transaction, so you can send the same many
// times in the same block, without it having the same Tx hash.
type Nonce struct {
	Value string `json:"value"`
}

// belongs to `system`  structs
type EOSIOParameters struct {
	TargetBlockSize              uint32 `json:"target_block_size"`
	MaxBlockSize                 uint32 `json:"max_block_size"`
	TargetBlockActsPerScope      uint32 `json:"target_block_acts_per_scope"`
	MaxBlockActsPerScope         uint32 `json:"max_block_acts_per_scope"`
	TargetBlockActs              uint32 `json:"target_block_acts"`
	MaxBlockActs                 uint32 `json:"max_block_acts"`
	MaxStorageSize               uint64 `json:"max_storage_size"`
	MaxTransactionLifetime       uint32 `json:"max_transaction_lifetime"`
	MaxTransactionExecTime       uint32 `json:"max_transaction_exec_time"`
	MaxAuthorityDepth            uint16 `json:"max_authority_depth"`
	MaxInlineDepth               uint16 `json:"max_inline_depth"`
	MaxInlineActionSize          uint32 `json:"max_inline_action_size"`
	MaxGeneratedTransactionSize  uint32 `json:"max_generated_transaction_size"`
	MaxGeneratedTransactionCount uint32 `json:"max_generated_transaction_count"`
	PercentOfMaxInflationRate    uint32 `json:"percent_of_max_inflation_rate"`
	StorageReserveRatio          uint32 `json:"storage_reserve_ratio"`
}

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

//
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

func init() {
	RegisterAction(AccountName("eosio"), ActionName("transfer"), &Transfer{})
	RegisterAction(AccountName("eosio"), ActionName("issue"), &Issue{})
	RegisterAction(AccountName("eosio"), ActionName("setcode"), &SetCode{})
	RegisterAction(AccountName("eosio"), ActionName("setabi"), &SetABI{})
	RegisterAction(AccountName("eosio"), ActionName("newaccount"), &SetABI{})
	RegisterAction(AccountName("eosio"), ActionName("delegatebw"), &DelegateBW{})
	RegisterAction(AccountName("eosio"), ActionName("undelegatebw"), &UndelegateBW{})
	RegisterAction(AccountName("eosio"), ActionName("refund"), &Refund{})
	RegisterAction(AccountName("eosio"), ActionName("regproducer"), &RegisterProducer{})
	RegisterAction(AccountName("eosio"), ActionName("unregprod"), &UnregisterProducer{})
	RegisterAction(AccountName("eosio"), ActionName("regproxy"), &RegisterProxy{})
	RegisterAction(AccountName("eosio"), ActionName("unregproxy"), &UnregisterProxy{})
	RegisterAction(AccountName("eosio"), ActionName("voteproducer"), &VoteProducer{})
	RegisterAction(AccountName("eosio"), ActionName("claimrewards"), &ClaimRewards{})
	RegisterAction(AccountName("eosio"), ActionName("nonce"), &Nonce{})
}

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
