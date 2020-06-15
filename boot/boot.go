package boot

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type option func(b *Boot) *Boot


func WithKeyBag(keyBag *eos.KeyBag) option {
	return func(b *Boot) *Boot {
		b.keyBag = keyBag
		return b
	}
}

func WithBootstrapping(genesisPath string, privateKey *nil) option {
	return func(b *Boot) *Boot {
		b.bootstrappingEnabled = true
		b.genesisPath = genesisPath
		return b
	}
}

func WithCachePath(cachePath string) option {
	return func(b *Boot) *Boot {
		b.cachePath = cachePath
		return b
	}
}

type Boot struct {
	bootSequencePath     string
	targetNetAPI         *eos.API
	bootstrappingEnabled bool
	genesisPath          string
	genesis              *GenesisJSON
	cachePath            string // Directory to store downloaded contract and ABI data
	bootSequence         *BootSeq

	keyBag *eos.KeyBag

	Snapshot           Snapshot
	WriteActions       bool
	HackVotingAccounts bool

	//Ephemeral
	//Ephemeral
}

func New(bootSequencePath string, targetAPI *eos.API, opts ...option) *Boot {
	b := &Boot{
		targetNetAPI:     targetAPI,
		bootSequencePath: bootSequencePath,
		cachePath:        "./boot-cache",
	}
	for _, opt := range opts {
		b = opt(b)
	}
	return b
}

func (b *Boot) getPublicKey() ecc.PublicKey {
	return b.keyBag.Keys[0].PublicKey()
}

func (b *Boot) getPrivateKey() *ecc.PrivateKey {
	return b.keyBag.Keys[0]
}

func (b *Boot) Run() (err error) {
	ctx := context.Background()

	b.bootSequence, err = ReadBootSeq(b.bootSequencePath)
	if err != nil {
		return err
	}

	zlog.Debug("downloading references")
	if err := b.downloadReferences(); err != nil {
		return err
	}

	zlog.Debug("setting boot keys")
	if err := b.setKeys(); err != nil {
		return err
	}

	if err := b.attachKeysOnTargetNode(ctx); err != nil {
		return err
	}

	if b.bootstrappingEnabled {
f		zlog.Debug("bootstrapping chain")
		if err := b.runBoostrapping(ctx); err != nil {
			return err
		}
	}

	b.pingTargetNetwork()

	zlog.Info("In-memory keys:")
	memkeys, _ := b.targetNetAPI.Signer.AvailableKeys(ctx)
	for _, key := range memkeys {
		zlog.Info("", zap.String("key", key.String()))
	}

	//eos.Debug = true

	for _, step := range b.bootSequence.BootSequence {
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
					_, err := b.targetNetAPI.SignPushActions(ctx, chunk...)
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

func (b *Boot) RunChainValidation() (bool, error) {
	bootSeqMap := ActionMap{}
	bootSeq := []*eos.Action{}

	for _, step := range b.bootSequence.BootSequence {
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

func (b *Boot) writeAllActionsToDisk() error {
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

	for _, step := range b.bootSequence.BootSequence {
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

func (b *Boot) pingTargetNetwork() {
	zlog.Info("Pinging target node at ", zap.String("url", b.targetNetAPI.BaseURL))
	for {
		info, err := b.targetNetAPI.GetInfo(context.Background())
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

func (b *Boot) validateTargetNetwork(bootSeqMap ActionMap, bootSeq []*eos.Action) (err error) {
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

		m, err := b.targetNetAPI.GetBlockByNum(context.Background(), uint32(blockHeight))
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

func (b *Boot) flushMissingActions(seenMap map[string]bool, bootSeq []*eos.Action) {
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

func (b *Boot) GetContentsCacheRef(filename string) (string, error) {
	for _, fl := range b.bootSequence.Contents {
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

func (b *Boot) writeToFile(filename, content string) {
	fl, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		zlog.Info("Unable to write to file", zap.String("file_name", filename), zap.Error(err))
		return
	}
	defer fl.Close()

	fl.Write([]byte(content))

	zlog.Info("Wrote file", zap.String("file_name", filename))
}
