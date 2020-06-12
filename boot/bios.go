package boot

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type BIOS struct {
	CachePath string

	TargetNetAPI       *eos.API
	Snapshot           Snapshot
	BootSequenceFile   string
	BootSequence       *BootSeq
	WriteActions       bool
	HackVotingAccounts bool
	ReuseGenesis       bool

	Genesis *GenesisJSON

	EphemeralPrivateKey *ecc.PrivateKey
	EphemeralPublicKey  ecc.PublicKey
}

func NewBIOS(cachePath string, targetAPI *eos.API) *BIOS {
	b := &BIOS{
		CachePath:    cachePath,
		TargetNetAPI: targetAPI,
	}
	return b
}

func (b *BIOS) Boot() error {
	bootSeq, err := ReadBootSeq(b.BootSequenceFile)
	if err != nil {
		return err
	}
	b.BootSequence = bootSeq

	if err := b.DownloadReferences(); err != nil {
		return err
	}

	zlog.Info("***************************************************************")
	zlog.Info("START BOOT SEQUENCE...")
	zlog.Info("***************************************************************")

	var genesisData string
	var pubKey ecc.PublicKey
	var privKey string

	err = b.setEphemeralKeypair()
	if err != nil {
		return err
	}

	ctx := context.Background()
	pubKey = b.EphemeralPublicKey
	privKey = b.EphemeralPrivateKey.String()

	if b.ReuseGenesis {
		genesisData, err = b.LoadGenesisFromFile(pubKey.String())
		if err != nil {
			return err
		}
	} else {
		genesisData = b.GenerateGenesisJSON(pubKey.String())

		b.writeToFile("genesis.pub", pubKey.String())
		b.writeToFile("genesis.key", privKey)
	}

	// Don't get `get_required_keys` from the blockchain, this adds
	// latency.. and we KNOW the key you're going to ask :) It's the
	// only key we're going to sign with anyway..
	b.TargetNetAPI.SetCustomGetRequiredKeys(func(ctx context.Context, tx *eos.Transaction) (out []ecc.PublicKey, err error) {
		return append(out, pubKey), nil
	})

	// Store keys in wallet, to sign `SetCode` and friends..
	if err := b.TargetNetAPI.Signer.ImportPrivateKey(ctx, privKey); err != nil {
		return fmt.Errorf("ImportWIF: %s", err)
	}

	if err := b.writeAllActionsToDisk(); err != nil {
		return fmt.Errorf("writing actions to disk: %s", err)
	}

	if err := b.DispatchBootNode(genesisData, pubKey.String(), privKey); err != nil {
		return fmt.Errorf("dispatch boot_node hook: %s", err)
	}

	b.pingTargetNetwork()

	zlog.Info("In-memory keys:")
	memkeys, _ := b.TargetNetAPI.Signer.AvailableKeys(ctx)
	for _, key := range memkeys {
		zlog.Info("", zap.String("key", key.String()))
	}

	//eos.Debug = true

	for _, step := range b.BootSequence.BootSequence {
		zlog.Info("action", zap.String("label", step.Label), zap.String("op", step.Op))

		acts, err := step.Data.Actions(b)
		if err != nil {
			return fmt.Errorf("getting actions for step %q: %s", step.Op, err)
		}

		if len(acts) != 0 {
			for idx, chunk := range ChunkifyActions(acts) {
				for _, c := range chunk {
					zlog.Info("processing chunk", zap.String("action", string(c.Name)))
				}
				err := Retry(25, time.Second, func() error {
					_, err := b.TargetNetAPI.SignPushActions(ctx, chunk...)
					if err != nil {
						zlog.Error("error pushing transaction", zap.String("op", step.Op), zap.Int("idx", idx), zap.Error(err))
						return fmt.Errorf("push actions for step %q, chunk %d: %s", step.Op, idx, err)
					}
					return nil
				})
				if err != nil {
					zlog.Info(" failed")
					return err
				}
			}
		}
	}

	zlog.Info("Waiting 2 seconds for transactions to flush to blocks")
	time.Sleep(2 * time.Second)

	// FIXME: don't do chain validation here..
	isValid, err := b.RunChainValidation()
	if err != nil {
		return fmt.Errorf("chain validation: %s", err)
	}
	if !isValid {
		zlog.Info("WARNING: chain invalid, destroying network if possible")
		os.Exit(0)
	}

	return nil
}

func (b *BIOS) setEphemeralKeypair() error {
	if b.EphemeralPrivateKey != nil {
		b.logEphemeralKey("Using preset key")
		return nil
	}

	if _, ok := b.BootSequence.Keys["ephemeral"]; ok {
		cnt := b.BootSequence.Keys["ephemeral"]
		privKey, err := ecc.NewPrivateKey(strings.TrimSpace(cnt))
		if err != nil {
			return fmt.Errorf("unable to correctly decode ephemeral private key %q: %s", cnt, err)
		}

		b.EphemeralPrivateKey = privKey
		b.EphemeralPublicKey = privKey.PublicKey()

		b.logEphemeralKey("Using user provider custom ephemeral keys from boot sequence")

	} else if b.ReuseGenesis {
		genesisPrivateKey, err := readPrivKeyFromFile("genesis.key")
		if err != nil {
			return err
		}

		b.EphemeralPrivateKey = genesisPrivateKey
		b.EphemeralPublicKey = genesisPrivateKey.PublicKey()

		b.logEphemeralKey("REUSING previously generated ephemeral keys from genesis")

	} else {
		ephemeralPrivateKey, err := b.GenerateEphemeralPrivKey()
		if err != nil {
			return err
		}

		b.EphemeralPrivateKey = ephemeralPrivateKey
		b.EphemeralPublicKey = ephemeralPrivateKey.PublicKey()

		b.logEphemeralKey("Generated ephemeral keys")
	}

	return nil
}

func (b *BIOS) logEphemeralKey(tag string) {
	pubKey := b.EphemeralPublicKey.String()
	privKey := b.EphemeralPrivateKey.String()

	zlog.Info("ephemeral key", zap.String("tag", tag), zap.String("pub_key", pubKey), zap.String("priv_key_prefix", privKey[:4]), zap.String("priv_key", privKey[len(privKey)-4:]))
}

func (b *BIOS) RunChainValidation() (bool, error) {
	bootSeqMap := ActionMap{}
	bootSeq := []*eos.Action{}

	for _, step := range b.BootSequence.BootSequence {
		acts, err := step.Data.Actions(b)
		if err != nil {
			return false, fmt.Errorf("validating: getting actions for step %q: %s", step.Op, err)
		}

		for _, stepAction := range acts {
			if stepAction == nil {
				continue
			}

			stepAction.SetToServer(true)
			data, err := eos.MarshalBinary(stepAction)
			if err != nil {
				return false, fmt.Errorf("validating: binary marshalling: %s", err)
			}
			key := sha2(data)

			// if _, ok := bootSeqMap[key]; ok {
			// 	// TODO: don't fatal here plz :)
			// 	log.Fatalf("Same action detected twice [%s] with key [%s]\n", stepAction.Name, key)
			// }
			bootSeqMap[key] = stepAction
			bootSeq = append(bootSeq, stepAction)
		}

	}

	err := b.validateTargetNetwork(bootSeqMap, bootSeq)
	if err != nil {
		zlog.Info("BOOT SEQUENCE VALIDATION FAILED:", zap.Error(err))
		return false, nil
	}

	zlog.Info("")
	zlog.Info("All good! Chain validation succeeded!")
	zlog.Info("")

	return true, nil
}

func (b *BIOS) writeAllActionsToDisk() error {
	if !b.WriteActions {
		zlog.Info("Not writing actions to 'actions.jsonl'. Activate with --write-actions")
		return nil
	}

	zlog.Info("Writing all actions to 'actions.jsonl'...")
	fl, err := os.Create("actions.jsonl")
	if err != nil {
		return err
	}
	defer fl.Close()

	for _, step := range b.BootSequence.BootSequence {
		acts, err := step.Data.Actions(b)
		if err != nil {
			return fmt.Errorf("fetch step %q: %s", step.Op, err)
		}

		for _, stepAction := range acts {
			if stepAction == nil {
				continue
			}

			stepAction.SetToServer(false)
			data, err := json.Marshal(stepAction)
			if err != nil {
				return fmt.Errorf("binary marshalling: %s", err)
			}

			_, err = fl.Write(data)
			if err != nil {
				return err
			}
			_, _ = fl.Write([]byte("\n"))
		}
	}

	return nil
}

type ActionMap map[string]*eos.Action

type ValidationError struct {
	Err               error
	BlockNumber       int
	Action            *eos.Action
	RawAction         []byte
	Index             int
	ActionHexData     string
	PackedTransaction *eos.PackedTransaction
}

func (e ValidationError) Error() string {
	s := fmt.Sprintf("Action [%d][%s::%s] absent from blocks\n", e.Index, e.Action.Account, e.Action.Name)

	data, err := json.Marshal(e.Action)
	if err != nil {
		s += fmt.Sprintf("    json generation err: %s\n", err)
	} else {
		s += fmt.Sprintf("    json data: %s\n", string(data))
	}
	s += fmt.Sprintf("    hex data: %s\n", hex.EncodeToString(e.RawAction))
	s += fmt.Sprintf("    error: %s\n", e.Err.Error())

	return s
}

type ValidationErrors struct {
	Errors []error
}

func (v ValidationErrors) Error() string {
	s := ""
	for _, err := range v.Errors {
		s += ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n"
		s += err.Error()
		s += "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n"
	}

	return s
}

func (b *BIOS) pingTargetNetwork() {
	zlog.Info("Pinging target node at ", zap.String("url", b.TargetNetAPI.BaseURL))
	for {
		info, err := b.TargetNetAPI.GetInfo(context.Background())
		if err != nil {
			zlog.Warn("target node", zap.Error(err))
			time.Sleep(1 * time.Second)
			continue
		}

		if info.HeadBlockNum < 2 {
			zlog.Info("target node: still no blocks in")
			zlog.Info(".")
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	zlog.Info(" touchdown!")
}

func (b *BIOS) validateTargetNetwork(bootSeqMap ActionMap, bootSeq []*eos.Action) (err error) {
	expectedActionCount := len(bootSeq)
	validationErrors := make([]error, 0)

	b.pingTargetNetwork()

	// TODO: wait for target network to be up, and responding...
	zlog.Info("Pulling blocks from chain until we gathered all actions to validate:")
	blockHeight := 1
	actionsRead := 0
	seenMap := map[string]bool{}
	gotSomeTx := false
	backOff := false
	timeBetweenFetch := time.Duration(0)
	var timeLastNotFound time.Time

	for {
		time.Sleep(timeBetweenFetch)

		m, err := b.TargetNetAPI.GetBlockByNum(context.Background(), uint32(blockHeight))
		if err != nil {
			if gotSomeTx && !backOff {
				backOff = true
				timeBetweenFetch = 500 * time.Millisecond
				timeLastNotFound = time.Now()

				time.Sleep(2000 * time.Millisecond)

				continue
			}

			zlog.Warn("Failed getting block num from target api", zap.String("message", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}

		blockHeight++

		zlog.Info("Receiving block", zap.Uint32("block_num", m.BlockNumber()), zap.String("producer", string(m.Producer)), zap.Int("trx_count", len(m.Transactions)))

		if !gotSomeTx && len(m.Transactions) > 2 {
			gotSomeTx = true
		}

		if !timeLastNotFound.IsZero() && timeLastNotFound.Before(time.Now().Add(-10*time.Second)) {
			b.flushMissingActions(seenMap, bootSeq)
		}

		for _, receipt := range m.Transactions {
			unpacked, err := receipt.Transaction.Packed.Unpack()
			if err != nil {
				zlog.Warn("Unable to unpack transaction, won't be able to fully validate", zap.Error(err))
				return fmt.Errorf("unpack transaction failed")
			}

			for _, act := range unpacked.Actions {
				act.SetToServer(false)
				data, err := eos.MarshalBinary(act)
				if err != nil {
					zlog.Error("Error marshalling an action", zap.Error(err))
					validationErrors = append(validationErrors, ValidationError{
						Err:               err,
						BlockNumber:       1, // extract from the block transactionmroot
						PackedTransaction: receipt.Transaction.Packed,
						Action:            act,
						RawAction:         data,
						ActionHexData:     hex.EncodeToString(act.HexData),
						Index:             actionsRead,
					})
					return err
				}
				key := sha2(data) // TODO: compute a hash here..

				if _, ok := bootSeqMap[key]; !ok {
					validationErrors = append(validationErrors, ValidationError{
						Err:               errors.New("not found"),
						BlockNumber:       1, // extract from the block transactionmroot
						PackedTransaction: receipt.Transaction.Packed,
						Action:            act,
						RawAction:         data,
						ActionHexData:     hex.EncodeToString(act.HexData),
						Index:             actionsRead,
					})
					zlog.Warn("INVALID action", zap.Int("action_read", actionsRead+1), zap.Int("expected_action_count", expectedActionCount), zap.String("account", string(act.Account)), zap.String("action", string(act.Name)))
				} else {
					seenMap[key] = true
					zlog.Info("validated action", zap.Int("action_read", actionsRead+1), zap.Int("expected_action_count", expectedActionCount), zap.String("account", string(act.Account)), zap.String("action", string(act.Name)))
				}

				actionsRead++
			}
		}

		if actionsRead == len(bootSeq) {
			break
		}

	}

	if len(validationErrors) > 0 {
		return ValidationErrors{Errors: validationErrors}
	}

	return nil
}

func (b *BIOS) flushMissingActions(seenMap map[string]bool, bootSeq []*eos.Action) {
	fl, err := os.Create("missing_actions.jsonl")
	if err != nil {
		zlog.Error("Couldn't write to `missing_actions.jsonl`:", zap.Error(err))
		return
	}
	defer fl.Close()

	// TODO: print all actions that are still MISSING to `missing_actions.jsonl`.
	zlog.Info("Flushing unseen transactions to `missing_actions.jsonl` up until this point.")

	for _, act := range bootSeq {
		act.SetToServer(true)
		data, _ := eos.MarshalBinary(act)
		key := sha2(data)

		if !seenMap[key] {
			act.SetToServer(false)
			data, _ := json.Marshal(act)
			fl.Write(data)
			fl.Write([]byte("\n"))
		}
	}
}

func (b *BIOS) inputGenesisData() (genesis *GenesisJSON) {
	zlog.Info("")
	for {
		zlog.Info("Please input the genesis data of the network you want to join: ")
		genesisData, err := ScanSingleLine()
		if err != nil {
			zlog.Error("error reading:", zap.Error(err))
			continue
		}

		err = json.Unmarshal([]byte(genesisData), &genesis)
		if err != nil {
			zlog.Error("Invalid genesis data", zap.Error(err))
			continue
		}

		return
	}
}

func (b *BIOS) GenerateEphemeralPrivKey() (*ecc.PrivateKey, error) {
	return ecc.NewRandomPrivateKey()
}

func (b *BIOS) GenerateGenesisJSON(pubKey string) string {
	// known not to fail
	cnt, _ := json.Marshal(&GenesisJSON{
		InitialTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05"),
		InitialKey:       pubKey,
	})
	return string(cnt)
}

func (b *BIOS) LoadGenesisFromFile(pubkey string) (string, error) {
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

func (b *BIOS) GetContentsCacheRef(filename string) (string, error) {
	for _, fl := range b.BootSequence.Contents {
		if fl.Name == filename {
			return fl.URL, nil
		}
	}
	return "", fmt.Errorf("%q not found in target contents", filename)
}

func ChunkifyActions(actions []*eos.Action) (out [][]*eos.Action) {
	currentChunk := []*eos.Action{}
	for _, act := range actions {
		if act == nil {
			if len(currentChunk) != 0 {
				out = append(out, currentChunk)
			}
			currentChunk = []*eos.Action{}
		} else {
			currentChunk = append(currentChunk, act)
		}
	}
	if len(currentChunk) > 0 {
		out = append(out, currentChunk)
	}
	return
}

func accountVariation(acct eos.AccountName, variation int) eos.AccountName {
	name := string(acct)
	if len(name) > 11 {
		name = name[:11]
	}
	variedName := name + string([]byte{'a' + byte(variation-1)})

	return eos.AccountName(variedName)
}

func readPrivKeyFromFile(filename string) (*ecc.PrivateKey, error) {
	cnt, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	strCnt := strings.TrimSpace(string(cnt))

	return ecc.NewPrivateKey(strCnt)
}

func (b *BIOS) writeToFile(filename, content string) {
	fl, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		zlog.Info("Unable to write to file", zap.String("file_name", filename), zap.Error(err))
		return
	}
	defer fl.Close()

	fl.Write([]byte(content))

	zlog.Info("Wrote file", zap.String("file_name", filename))
}
