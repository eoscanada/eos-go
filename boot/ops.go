package boot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"go.uber.org/zap"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/eoscanada/eos-go/token"
	"github.com/eoscanada/eosc/bios/unregd"
)

type Operation interface {
	Actions(b *Boot) ([]*eos.Action, error)
}

var operationsRegistry = map[string]Operation{
	"system.setcode":             &OpSetCode{},
	"system.setram":              &OpSetRAM{},
	"system.newaccount":          &OpNewAccount{},
	"system.setpriv":             &OpSetPriv{},
	"token.create":               &OpCreateToken{},
	"token.issue":                &OpIssueToken{},
	"token.transfer":             &OpTransferToken{},
	"system.setprods":            &OpSetProds{},
	"snapshot.create_accounts":   &OpSnapshotCreateAccounts{},
	"snapshot.load_unregistered": &OpInjectUnregdSnapshot{},
	"system.resign_accounts":     &OpResignAccounts{},
	"system.create_voters":       &OpCreateVoters{},
	"system.delegate_bw":         &OpDelegateBW{},
	"system.buy_ram":             &OpBuyRam{},
	"system.buy_ram_bytes":       &OpBuyRamBytes{},
}

type OperationType struct {
	Op    string
	Label string
	Data  Operation
}

func (o *OperationType) UnmarshalJSON(data []byte) error {
	opData := struct {
		Op    string
		Label string
		Data  json.RawMessage
	}{}
	if err := json.Unmarshal(data, &opData); err != nil {
		return err
	}

	opType, found := operationsRegistry[opData.Op]
	if !found {
		return fmt.Errorf("operation type %q invalid, use one of: %q", opData.Op, operationsRegistry)
	}

	objType := reflect.TypeOf(opType).Elem()
	obj := reflect.New(objType).Interface()

	if len(opData.Data) != 0 {
		err := json.Unmarshal(opData.Data, &obj)
		if err != nil {
			return fmt.Errorf("operation type %q invalid, error decoding: %s", opData.Op, err)
		}
	} //  else {
	// 	_ = json.Unmarshal([]byte("{}"), &obj)
	// }

	opIface, ok := obj.(Operation)
	if !ok {
		return fmt.Errorf("operation type %q isn't an op", opData.Op)
	}

	*o = OperationType{
		Op:    opData.Op,
		Label: opData.Label,
		Data:  opIface,
	}

	return nil
}

//

type OpSetCode struct {
	Account         eos.AccountName
	ContractNameRef string `json:"contract_name_ref"`
}

func (op *OpSetCode) Actions(b *Boot) (out []*eos.Action, err error) {
	wasmFileRef, err := b.GetContentsCacheRef(fmt.Sprintf("%s.wasm", op.ContractNameRef))
	if err != nil {
		return nil, err
	}
	abiFileRef, err := b.GetContentsCacheRef(fmt.Sprintf("%s.abi", op.ContractNameRef))
	if err != nil {
		return nil, err
	}

	setCode, err := system.NewSetCodeTx(
		op.Account,
		b.FileNameFromCache(wasmFileRef),
		b.FileNameFromCache(abiFileRef),
	)
	if err != nil {
		return nil, fmt.Errorf("NewSetCodeTx %s: %s", op.ContractNameRef, err)
	}

	return setCode.Actions, nil
}

//

type OpSetRAM struct {
	MaxRAMSize uint64 `json:"max_ram_size"`
}

func (op *OpSetRAM) Actions(b *Boot) (out []*eos.Action, err error) {
	return append(out, system.NewSetRAM(op.MaxRAMSize)), nil
}

//

type OpNewAccount struct {
	Creator    eos.AccountName
	NewAccount eos.AccountName `json:"new_account"`
	Pubkey     string
	RamBytes   uint32 `json:"ram_bytes"`
}

func (op *OpNewAccount) Actions(b *Boot) (out []*eos.Action, err error) {
	pubKey := b.getPublicKey()

	if op.Pubkey != "ephemeral" {
		pubKey, err = ecc.NewPublicKey(op.Pubkey)
		if err != nil {
			return nil, fmt.Errorf("reading pubkey: %s", err)
		}
	}
	out = append(out, system.NewNewAccount(op.Creator, op.NewAccount, pubKey))

	if op.RamBytes > 0 {
		out = append(out, system.NewBuyRAMBytes(op.Creator, op.NewAccount, op.RamBytes))
	}

	return out, nil
}

type OpDelegateBW struct {
	From     eos.AccountName
	To       eos.AccountName
	StakeCPU int64 `json:"stake_cpu"`
	StakeNet int64 `json:"stake_net"`
	Transfer bool
}

func (op *OpDelegateBW) Actions(b *Boot) (out []*eos.Action, err error) {
	return append(out, system.NewDelegateBW(op.From, op.To, eos.NewEOSAsset(op.StakeCPU), eos.NewEOSAsset(op.StakeNet), op.Transfer)), nil
}

type OpBuyRam struct {
	Payer       eos.AccountName
	Receiver    eos.AccountName
	EOSQuantity uint64 `json:"eos_quantity"`
}

func (op *OpBuyRam) Actions(b *Boot) (out []*eos.Action, err error) {
	return append(out, system.NewBuyRAM(op.Payer, op.Receiver, op.EOSQuantity)), nil
}

type OpBuyRamBytes struct {
	Payer    eos.AccountName
	Receiver eos.AccountName
	Bytes    uint32
}

func (op *OpBuyRamBytes) Actions(b *Boot) (out []*eos.Action, err error) {
	return append(out, system.NewBuyRAMBytes(op.Payer, op.Receiver, op.Bytes)), nil
}

type OpTransfer struct {
	From   eos.AccountName
	to     eos.AccountName
	Amount eos.Asset
	Memo   string
}

func (op *OpTransfer) Actions(b *Boot) (out []*eos.Action, err error) {
	return append(out, token.NewTransfer(op.From, op.to, op.Amount, op.Memo)), nil
}

type OpCreateVoters struct {
	Creator eos.AccountName
	Pubkey  string
	Count   int
}

func (op *OpCreateVoters) Actions(b *Boot) (out []*eos.Action, err error) {
	pubKey := b.getPublicKey()

	if op.Pubkey != "ephemeral" {
		pubKey, err = ecc.NewPublicKey(op.Pubkey)
		if err != nil {
			return nil, fmt.Errorf("reading pubkey: %s", err)
		}
	}

	for i := 0; i < op.Count; i++ {
		voterName := eos.AccountName(voterName(i))
		fmt.Println("Creating voter: ", voterName)
		out = append(out, system.NewNewAccount(op.Creator, voterName, pubKey))
		out = append(out, token.NewTransfer(op.Creator, voterName, eos.NewEOSAsset(1000000000), ""))
		out = append(out, system.NewBuyRAMBytes(AN("eosio"), voterName, 8192)) // 8kb gift !
		out = append(out, system.NewDelegateBW(AN("eosio"), voterName, eos.NewEOSAsset(10000), eos.NewEOSAsset(10000), true))
	}
	return
}

const charset = "abcdefghijklmnopqrstuvwxyz"

func voterName(index int) string {
	padding := string(bytes.Repeat([]byte{charset[index]}, 7))
	return "voter" + padding
}

type OpSetPriv struct {
	Account eos.AccountName
}

func (op *OpSetPriv) Actions(b *Boot) (out []*eos.Action, err error) {
	return append(out, system.NewSetPriv(op.Account)), nil
}

type OpCreateToken struct {
	Account eos.AccountName `json:"account"`
	Amount  eos.Asset       `json:"amount"`
}

func (op *OpCreateToken) Actions(b *Boot) (out []*eos.Action, err error) {
	act := token.NewCreate(op.Account, op.Amount)
	return append(out, act), nil
}

type OpIssueToken struct {
	Account eos.AccountName
	Amount  eos.Asset
	Memo    string
}

func (op *OpIssueToken) Actions(b *Boot) (out []*eos.Action, err error) {
	act := token.NewIssue(op.Account, op.Amount, op.Memo)
	return append(out, act), nil
}

//

type OpTransferToken struct {
	From     eos.AccountName
	To       eos.AccountName
	Quantity eos.Asset
	Memo     string
}

func (op *OpTransferToken) Actions(b *Boot) (out []*eos.Action, err error) {
	act := token.NewTransfer(op.From, op.To, op.Quantity, op.Memo)
	return append(out, act), nil
}

//

type OpSnapshotCreateAccounts struct {
	BuyRAMBytes             uint64 `json:"buy_ram_bytes"`
	TestnetTruncateSnapshot int    `json:"TESTNET_TRUNCATE_SNAPSHOT"`
}

func (op *OpSnapshotCreateAccounts) Actions(b *Boot) (out []*eos.Action, err error) {
	snapshotFile, err := b.GetContentsCacheRef("snapshot.csv")
	if err != nil {
		return nil, err
	}

	rawSnapshot, err := b.ReadFromCache(snapshotFile)
	if err != nil {
		return nil, fmt.Errorf("reading snapshot file: %s", err)
	}

	snapshotData, err := NewSnapshot(rawSnapshot)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot csv: %s", err)
	}

	if len(snapshotData) == 0 {
		return nil, fmt.Errorf("snapshot is empty or not loaded")
	}

	wellKnownPubkey, _ := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")

	for idx, hodler := range snapshotData {
		if trunc := op.TestnetTruncateSnapshot; trunc != 0 {
			if idx == trunc {
				zlog.Debug("truncated snapshot", zap.Int("at_row", trunc))
				break
			}
		}

		destAccount := AN(hodler.AccountName)
		destPubKey := hodler.EOSPublicKey
		if b.HackVotingAccounts {
			destPubKey = wellKnownPubkey
		}

		out = append(out, system.NewNewAccount(AN("eosio"), destAccount, destPubKey))

		cpuStake, netStake, rest := splitSnapshotStakes(hodler.Balance)

		// special case `transfer` for `b1` ?
		out = append(out, system.NewDelegateBW(AN("eosio"), destAccount, cpuStake, netStake, true))
		out = append(out, system.NewBuyRAMBytes(AN("eosio"), destAccount, uint32(op.BuyRAMBytes)))
		out = append(out, nil) // end transaction

		memo := "Welcome " + hodler.EthereumAddress[len(hodler.EthereumAddress)-6:]
		out = append(out, token.NewTransfer(AN("eosio"), destAccount, rest, memo), nil)
	}

	return
}

func splitSnapshotStakes(balance eos.Asset) (cpu, net, xfer eos.Asset) {
	if balance.Amount < 5000 {
		return
	}

	// everyone has minimum 0.25 EOS staked
	// some 10 EOS unstaked
	// the rest split between the two

	cpu = eos.NewEOSAsset(2500)
	net = eos.NewEOSAsset(2500)

	remainder := eos.NewEOSAsset(int64(balance.Amount - cpu.Amount - net.Amount))

	if remainder.Amount <= 100000 /* 10.0 EOS */ {
		return cpu, net, remainder
	}

	remainder.Amount -= 100000 // keep them floating, unstaked

	firstHalf := remainder.Amount / 2
	cpu.Amount += firstHalf
	net.Amount += remainder.Amount - firstHalf

	return cpu, net, eos.NewEOSAsset(100000)
}

//

type OpInjectUnregdSnapshot struct {
	TestnetTruncateSnapshot int `json:"TESTNET_TRUNCATE_SNAPSHOT"`
}

func (op *OpInjectUnregdSnapshot) Actions(b *Boot) (out []*eos.Action, err error) {
	snapshotFile, err := b.GetContentsCacheRef("snapshot_unregistered.csv")
	if err != nil {
		return nil, err
	}

	rawSnapshot, err := b.ReadFromCache(snapshotFile)
	if err != nil {
		return nil, fmt.Errorf("reading snapshot file: %s", err)
	}

	snapshotData, err := NewUnregdSnapshot(rawSnapshot)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot csv: %s", err)
	}

	if len(snapshotData) == 0 {
		return nil, fmt.Errorf("snapshot is empty or not loaded")
	}

	for idx, hodler := range snapshotData {
		if trunc := op.TestnetTruncateSnapshot; trunc != 0 {
			if idx == trunc {
				zlog.Debug("- DEBUG: truncated unreg'd snapshot", zap.Int("row", trunc))
				break
			}
		}

		//system.NewDelegatedNewAccount(AN("eosio"), AN(hodler.AccountName), AN("eosio.unregd"))

		out = append(out,
			unregd.NewAdd(hodler.EthereumAddress, hodler.Balance),
			token.NewTransfer(AN("eosio"), AN("eosio.unregd"), hodler.Balance, "Future claim"),
			nil,
		)
	}

	return
}

//

type producerKeyString struct {
	ProducerName          eos.AccountName `json:"producer_name"`
	BlockSigningKeyString string          `json:"block_signing_key"`
}

type OpSetProds struct {
	Prods []producerKeyString
}

func (op *OpSetProds) Actions(b *Boot) (out []*eos.Action, err error) {

	var prodKeys []system.ProducerKey

	for _, key := range op.Prods {
		prodKey := system.ProducerKey{
			ProducerName: key.ProducerName,
		}
		if key.BlockSigningKeyString == "" || key.BlockSigningKeyString == "ephemeral" {
			prodKey.BlockSigningKey = b.getPublicKey()
		} else {
			k, err := ecc.NewPublicKey(key.BlockSigningKeyString)
			if err != nil {
				panic(err)
			}
			prodKey.BlockSigningKey = k
		}
		prodKeys = append(prodKeys, prodKey)
	}

	if len(prodKeys) == 0 {
		prodKeys = []system.ProducerKey{system.ProducerKey{
			ProducerName:    AN("eosio"),
			BlockSigningKey: b.getPublicKey(),
		}}
	}

	var producers []string
	for _, key := range prodKeys {
		producers = append(producers, string(key.ProducerName))
	}
	fmt.Printf("Producers set to: %v\n", producers)

	out = append(out, system.NewSetProds(prodKeys))
	return
}

//

type OpResignAccounts struct {
	Accounts            []eos.AccountName
	TestnetKeepAccounts bool `json:"TESTNET_KEEP_ACCOUNTS"`
}

func (op *OpResignAccounts) Actions(b *Boot) (out []*eos.Action, err error) {
	if op.TestnetKeepAccounts {
		zlog.Debug("Keeping system accounts around, for testing purposes.")
		return
	}

	systemAccount := AN("eosio")
	prodsAccount := AN("eosio.prods") // this is a special system account that is granted by 2/3 + 1 of the current BP schedule.

	eosioPresent := false
	for _, acct := range op.Accounts {
		if acct == systemAccount {
			eosioPresent = true
			continue
		}

		out = append(out,
			system.NewUpdateAuth(acct, PN("active"), PN("owner"), eos.Authority{
				Threshold: 1,
				Accounts: []eos.PermissionLevelWeight{
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      AN("eosio"),
							Permission: PN("active"),
						},
						Weight: 1,
					},
				},
			}, PN("active")),
			system.NewUpdateAuth(acct, PN("owner"), PN(""), eos.Authority{
				Threshold: 1,
				Accounts: []eos.PermissionLevelWeight{
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      AN("eosio"),
							Permission: PN("active"),
						},
						Weight: 1,
					},
				},
			}, PN("owner")),
		)
	}

	if eosioPresent {
		out = append(out,
			system.NewUpdateAuth(systemAccount, PN("active"), PN("owner"), eos.Authority{
				Threshold: 1,
				Accounts: []eos.PermissionLevelWeight{
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      prodsAccount,
							Permission: PN("active"),
						},
						Weight: 1,
					},
				},
			}, PN("active")),
			system.NewUpdateAuth(systemAccount, PN("owner"), PN(""), eos.Authority{
				Threshold: 1,
				Accounts: []eos.PermissionLevelWeight{
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      prodsAccount,
							Permission: PN("active"),
						},
						Weight: 1,
					},
				},
			}, PN("owner")),
		)
	}

	out = append(out, nil)

	return
}
