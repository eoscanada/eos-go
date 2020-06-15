package boot

import (
	"context"
	"fmt"
	"strings"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"go.uber.org/zap"
)

func (b *Boot) setKeys() error {
	if b.keyBag != nil {
		b.logKey("Using preset key bag")
		return nil
	}

	if _, ok := b.bootSequence.Keys["boot"]; ok {
		cnt := b.bootSequence.Keys["boot"]
		privKey, err := ecc.NewPrivateKey(strings.TrimSpace(cnt))
		if err != nil {
			return fmt.Errorf("unable to correctly decode boot private key %q: %s", cnt, err)
		}

		b.keyBag = eos.NewKeyBag()
		b.keyBag.Add(privKey.String())
		b.logKey("Using user provider custom boot key from boot sequence")
	} else {
		return fmt.Errorf("no key specified, either specify a key within the bootsequence or use `withKeys` options when creating a new `Boot`")
	}

	return nil
}

func (b *Boot) attachKeysOnTargetNode(ctx context.Context) error {

	// Don't get `get_required_keys` from the blockchain, this adds
	// latency.. and we KNOW the key you're going to ask :) It's the
	// only key we're going to sign with anyway..
	b.targetNetAPI.SetCustomGetRequiredKeys(func(ctx context.Context, tx *eos.Transaction) (out []ecc.PublicKey, err error) {
		for _, k := range b.keyBag.Keys {
			out = append(out, k.PublicKey())
		}
		return out, nil
	})

	// Store keys in wallet, to sign `SetCode` and friends..
	b.targetNetAPI.SetSigner(b.keyBag)
	//for _, key := range b.keyBag.Keys {
	//	zlog.Info("importing private keys", zap.String("key", key.String()))
	//	if err := b.targetNetAPI.Signer.ImportPrivateKey(ctx, key.String()); err != nil {
	//		return fmt.Errorf("ImportWIF: %s", err)
	//	}
	//}
	return nil
}

func (b *Boot) logKey(tag string) {
	pubKey := b.getPublicKey().String()
	privKey := b.getPrivateKey().String()

	zlog.Info("boot key",
		zap.String("tag", tag),
		zap.String("pub_key", pubKey),
		zap.String("priv_key_prefix", privKey[:4]),
		zap.String("priv_key", privKey[len(privKey)-4:]),
	)
}
