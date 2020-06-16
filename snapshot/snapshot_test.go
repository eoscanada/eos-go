package snapshot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/eoscanada/eos-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshotRead(t *testing.T) {
	// "/tmp/0125111385-07750c59b24ed52d2dbf2048b67b58e9c9bd53ff5cc4550277718c1d5d800f73-snapshot.bin"
	readSnapshotFile := os.Getenv("READ_SNAPSHOT_FILE")
	if readSnapshotFile == "" || !fileExists(readSnapshotFile) {
		t.Skipf("Environment varaible 'READ_SNAPSHOT_FILE' not set or value %q is not an exisiting file", readSnapshotFile)
		return
	}
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "name",
			input:  "input",
			expect: "expect",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := NewReader(readSnapshotFile)
			defer r.Close()

			assert.NoError(t, err)
			assert.Equal(t, r.Header.Version, uint32(1))

			for {
				section, err := r.Next()
				if err == io.EOF {
					break
				}
				assert.NoError(t, err)
				fmt.Println("Section", section.Name, "rows", section.RowCount, "bytes", section.BufferSize)

				switch section.Name {
				case "eosio::chain::chain_snapshot_header":
				case "eosio::chain::block_state":
					cnt := make([]byte, section.BufferSize)
					_, err := section.Buffer.Read(cnt)
					//fmt.Println(hex.EncodeToString(cnt))
					require.NoError(t, err)
					var state eos.BlockState
					assert.NoError(t, eos.UnmarshalBinary(cnt, &state))
					cnt, _ = json.MarshalIndent(state, "  ", "  ")
					fmt.Println(string(cnt))

				case "eosio::chain::account_object":
					require.NoError(t, readAccountObjects(section.Buffer, section.RowCount))
					cnt := make([]byte, section.BufferSize)
					_, err := section.Buffer.Read(cnt)
					require.NoError(t, err)

					require.NoError(t, ioutil.WriteFile("/tmp/test.dat", cnt, 0664))

					var accounts []AccountObject
					assert.NoError(t, eos.UnmarshalBinary(cnt, &accounts))
					cnt, _ = json.MarshalIndent(accounts, "  ", "  ")
					fmt.Println(string(cnt))

				case "eosio::chain::account_metadata_object":
				case "eosio::chain::account_ram_correction_object":
				case "eosio::chain::global_property_object":
				case "eosio::chain::protocol_state_object":
				case "eosio::chain::dynamic_global_property_object":
				case "eosio::chain::block_summary_object":
				case "eosio::chain::transaction_object":
				case "eosio::chain::generated_transaction_object":
				case "eosio::chain::code_object":
				case "contract_tables":
					err := readContractTables(section.Buffer)
					require.NoError(t, err)
				case "eosio::chain::permission_object":
				case "eosio::chain::permission_link_object":
				case "eosio::chain::resource_limits::resource_limits_object":
				case "eosio::chain::resource_limits::resource_usage_object":
				case "eosio::chain::resource_limits::resource_limits_state_object":
				case "eosio::chain::resource_limits::resource_limits_config_object":
				default:
					panic("unsupported section")
				}
			}
			fmt.Println("Done")
		})
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		return false
	}

	return !info.IsDir()
}
