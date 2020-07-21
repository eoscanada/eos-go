package snapshot

import "io"

type Section struct {
	Name       string
	Offset     uint64
	Size       uint64 // This includes the section name and row count
	BufferSize uint64 // This represents the bytes that are following the section header
	RowCount   uint64 // This is a count of rows packed in `Buffer`
	Buffer     io.Reader
}

// Next reads the next row
func (s *Section) Next() ([]byte, error) {
	return nil, nil
}

type callbackFunc func(obj interface{}) error

func (s *Section) Process(f callbackFunc) error {
	switch s.Name {
	case "eosio::chain::chain_snapshot_header":
		return s.readChainSnapshotHeader(f)
	case "eosio::chain::block_state":
		return s.readBlockState(f)
	case "eosio::chain::account_object":
		return s.readAccountObjects(f)
	case "eosio::chain::account_metadata_object":
		return s.readAccountMetadataObjects(f)
	case "eosio::chain::account_ram_correction_object":
		return s.readAccountRAMCorrectionObject(f)
	case "eosio::chain::global_property_object":
		return s.readGlobalPropertyObject(f)
	case "eosio::chain::protocol_state_object":
		return s.readProtocolStateObject(f)
	case "eosio::chain::dynamic_global_property_object":
		return s.readDynamicGlobalPropertyObject(f)
	case "eosio::chain::block_summary_object":
		return s.readBlockSummary(f)
	case "eosio::chain::transaction_object":
		return s.readTransactionObject(f)
	case "eosio::chain::generated_transaction_object":
		return s.readGeneratedTransactionObject(f)
	case "eosio::chain::code_object":
		return s.readCodeObject(f)
	case "contract_tables":
		return s.readContractTables(f)
	case "eosio::chain::permission_object":
		return s.readPermissionObject(f)
	case "eosio::chain::permission_link_object":
		return s.readPermissionLinkObject(f)
	case "eosio::chain::resource_limits::resource_limits_object":
		return s.readResourceLimitsObject(f)
	case "eosio::chain::resource_limits::resource_usage_object":
		return s.readResourceUsageObject(f)
	case "eosio::chain::resource_limits::resource_limits_state_object":
		return s.readResourceLimitsStateObject(f)
	case "eosio::chain::resource_limits::resource_limits_config_object":
		return s.readResourceLimitsConfigObject(f)

	case "eosio::chain::genesis_state":
		return s.readGenesisState(f)

	default:
		panic("unsupported section: " + s.Name)
	}

}
