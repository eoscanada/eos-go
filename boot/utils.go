package boot

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	eos "github.com/eoscanada/eos-go"
	"gopkg.in/olivere/elastic.v3/backoff"
)

// AN is a shortcut to create an AccountName
var AN = eos.AN

// PN is a shortcut to create a PermissionName
var PN = eos.PN

func Retry(attempts int, sleep time.Duration, callback func() error) (err error) {
	b := backoff.NewExponentialBackoff(sleep, 5*time.Second)
	for i := 0; ; i++ {
		err = callback()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(b.Next())

		log.Println("retrying after error:", err)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func sha2(input []byte) string {
	hash := sha256.New()
	_, _ = hash.Write(input) // can't fail
	return hex.EncodeToString(hash.Sum(nil))
}
