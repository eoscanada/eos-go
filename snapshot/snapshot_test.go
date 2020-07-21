package snapshot

import (
	"fmt"
	"io"
	"os"
	"testing"

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
			filename := "/tmp/0125111385-07750c59b24ed52d2dbf2048b67b58e9c9bd53ff5cc4550277718c1d5d800f73-snapshot.bin" // mainnet
			//filename := "/tmp/0003212331-0031042b02b2cf711fee6e1e24da94101fa6c1ea9ece568d5f13232473429db1-snapshot.bin" // kylin
			r, err := NewReader(filename)
			fmt.Println("Filename", filename)
			defer r.Close()

			assert.NoError(t, err)
			assert.Equal(t, r.Header.Version, uint32(1))

			var totalsize = 0

			for {
				section, err := r.Next()
				if err == io.EOF {
					break
				}
				assert.NoError(t, err)
				fmt.Println("Section", section.Name, "rows", section.RowCount, "bytes", section.BufferSize, "offset", section.Offset)
				totalsize += int(section.BufferSize)

				switch section.Name {
				case "eosio::chain::chain_snapshot_header":
					require.NoError(t, readChainSnapshotHeader(section))
				case "eosio::chain::block_state":
					// require.NoError(t, readBlockState(section))
				case "eosio::chain::account_object":
					// require.NoError(t, readAccountObjects(section))
				case "eosio::chain::account_metadata_object":
					//require.NoError(t, readAccountMetadataObjects(section))
				case "eosio::chain::account_ram_correction_object":
					//require.NoError(t, readAccountRAMCorrectionObject(section))
				case "eosio::chain::global_property_object":
					//require.NoError(t, readGlobalPropertyObject(section))
				case "eosio::chain::protocol_state_object":
					//require.NoError(t, readProtocolStateObject(section))
				case "eosio::chain::dynamic_global_property_object":
					// require.NoError(t, readDynamicGlobalPropertyObject(section))
				case "eosio::chain::block_summary_object":
					//require.NoError(t, readBlockSummary(section))
				case "eosio::chain::transaction_object":
					require.NoError(t, readTransactionObject(section))
				case "eosio::chain::generated_transaction_object":
					// require.NoError(t, readGeneratedTransactionObject(section))
				case "eosio::chain::code_object":
					// require.NoError(t, readCodeObject(section))
				case "contract_tables":
					// require.NoError(t, readContractTables(section))
				case "eosio::chain::permission_object":
					// require.NoError(t, readPermissionObject(section))
				case "eosio::chain::permission_link_object":
					// require.NoError(t, readPermissionLinkObject(section))
				case "eosio::chain::resource_limits::resource_limits_object":
					// require.NoError(t, readResourceLimitsObject(section))
				case "eosio::chain::resource_limits::resource_usage_object":
					// require.NoError(t, readResourceUsageObject(section))
				case "eosio::chain::resource_limits::resource_limits_state_object":
					// require.NoError(t, readResourceLimitsStateObject(section))
				case "eosio::chain::resource_limits::resource_limits_config_object":
					// require.NoError(t, readResourceLimitsConfigObject(section))
				case "eosio::chain::genesis_state":
					// // THIS SEEMS TO EXIST ONLY IN VERSION 2 OF THE SNAPSHOT FILE FORMAT.
					// // FOR NOW, WE ARE CONCENTRATING ON VERSION 3 (latest)
					// cnt := make([]byte, section.BufferSize)
					// _, err := section.Buffer.Read(cnt)
					// require.NoError(t, err)

					// var state GenesisState
					// assert.NoError(t, eos.UnmarshalBinary(cnt, &state))
					// cnt, _ = json.MarshalIndent(state, "  ", "  ")
					// fmt.Println(string(cnt))

				default:
					panic("unsupported section: " + section.Name)
				}
			}
			fmt.Println("Done", totalsize)
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
