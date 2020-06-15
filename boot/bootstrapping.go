package boot

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/eoscanada/eos-go/ecc"

	"github.com/eoscanada/eos-go"
)

func (b *Boot) runBoostrapping(ctx context.Context) (err error) {
	var genesisData string

	// using the first key in the bad
	privateBootstrappingKey := b.getPrivateKey()

	if b.genesisPath != "" {
		genesisData, err = loadGenesisFromFile(privateBootstrappingKey.PublicKey().String())
		if err != nil {
			return err
		}
	} else {
		genesisData = generateGenesisJSON(b.getPublicKey().String())
		b.writeToFile("genesis.pub", privateBootstrappingKey.PublicKey().String())
		b.writeToFile("genesis.key", privateBootstrappingKey.String())
	}

	// Don't get `get_required_keys` from the blockchain, this adds
	// latency.. and we KNOW the key you're going to ask :) It's the
	// only key we're going to sign with anyway..
	b.targetNetAPI.SetCustomGetRequiredKeys(func(ctx context.Context, tx *eos.Transaction) (out []ecc.PublicKey, err error) {
		return append(out, b.getPublicKey()), nil
	})

	// Store keys in wallet, to sign `SetCode` and friends..
	for _, key := range b.keyBag.Keys {
		zlog.Info("importing private keys", zap.String("key", key.String()))
		if err := b.targetNetAPI.Signer.ImportPrivateKey(ctx, key.String()); err != nil {
			return fmt.Errorf("ImportWIF: %s", err)
		}
	}

	if err := b.DispatchBootNode(genesisData, b.getPublicKey().String(), b.getPrivateKey().String()); err != nil {
		return fmt.Errorf("dispatch boot_node hook: %s", err)
	}

	return nil

}
