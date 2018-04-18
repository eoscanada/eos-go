package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

// A simple Wallet implementation, using a KeyBag and the `btcec`
// signature machinery.

func main() {

	keyBag := eos.NewKeyBag()
	for _, key := range []string{
		"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3",
		"5K7Ffo8LXHhbsxV48w3sZzo8UnaKX3z5iD5mvac1AfDhHXKs3ao",
		"5Jv2xEfJ4UVbNTBNjMxdWZAJcaDweP4bgwRd554NpWG3VynxW6L", // EOS5Dg9cu3yn5cMpWkrZnhmYk2xDBWmu62Sj2dNrWn6Ui82eoYJQh
		"5KE5hGNCAs1YvV74Ho14y1rV1DrnqZpTwLugS8QvYbKbrGAvVA1", // EOS71W8hvF43Eq6GQBRhuc5mvWKtknxzmb9NzNwPGpcEm2xAZaG8c
	} {
		if err := keyBag.Add(key); err != nil {
			log.Fatalln("Couldn't load private key:", err)
		}
	}

	http.HandleFunc("/v1/wallet/get_public_keys", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling get_public_keys")

		var out []string
		for _, key := range keyBag.Keys {
			out = append(out, key.PublicKey().String())
		}
		json.NewEncoder(w).Encode(out)
	})

	http.HandleFunc("/v1/wallet/sign_transaction", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling sign_transaction")

		var inputs []json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
			fmt.Println("sign_transaction: error:", err)
			http.Error(w, "couldn't decode input", 500)
			return
		}

		var tx *eos.SignedTransaction
		var requiredKeys []ecc.PublicKey
		var chainID eos.HexBytes

		if len(inputs) != 3 {
			http.Error(w, "invalid length of message, should be 3 parameters", 500)
			return
		}

		err := json.Unmarshal(inputs[0], &tx)
		if err != nil {
			http.Error(w, "decoding param1", 500)
			return
		}

		err = json.Unmarshal(inputs[1], &requiredKeys)
		if err != nil {
			http.Error(w, "decoding param1", 500)
			return
		}

		err = json.Unmarshal(inputs[2], &chainID)
		if err != nil {
			http.Error(w, "decoding param1", 500)
			return
		}

		signed, err := keyBag.Sign(tx, chainID, requiredKeys...)
		if err != nil {
			http.Error(w, fmt.Sprintf("error signing: %s", err), 500)
			return
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(signed)
	})

	http.HandleFunc("/v1/wallet/import_key", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling import_key")
		// TODO: do the things..
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)
		// We're ignoring inputs[0] which is the name of the wallet ("default" by default)
		keyBag.Add(inputs[1])
		w.WriteHeader(201)
		w.Write([]byte("{}"))
	})

	fmt.Println("Listening for wallet operations on port 5555")
	if err := http.ListenAndServe("127.0.0.1:5555", nil); err != nil {
		log.Println("Litsening failed:", err)
	}
}
