package snapshot

import "io"

type SectionName string

const (
	SectionNameChainSnapshotHeader         SectionName = "eosio::chain::chain_snapshot_header"
	SectionNameBlockState                  SectionName = "eosio::chain::block_state"
	SectionNameAccountObject               SectionName = "eosio::chain::account_object"
	SectionNameAccountMetadataObject       SectionName = "eosio::chain::account_metadata_object"
	SectionNameAccountRamCorrectionObject  SectionName = "eosio::chain::account_ram_correction_object"
	SectionNameGlobalPropertyObject        SectionName = "eosio::chain::global_property_object"
	SectionNameProtocolStateObject         SectionName = "eosio::chain::protocol_state_object"
	SectionNameDynamicGlobalPropertyObject SectionName = "eosio::chain::dynamic_global_property_object"
	SectionNameBlockSummaryObject          SectionName = "eosio::chain::block_summary_object"
	SectionNameTransactionObject           SectionName = "eosio::chain::transaction_object"
	SectionNameGeneratedTransactionObject  SectionName = "eosio::chain::generated_transaction_object"
	SectionNameCodeObject                  SectionName = "eosio::chain::code_object"
	SectionNameContractTables              SectionName = "contract_tables"
	SectionNamePermissionObject            SectionName = "eosio::chain::permission_object"
	SectionNamePermissionLinkObject        SectionName = "eosio::chain::permission_link_object"
	SectionNameResourceLimitsObject        SectionName = "eosio::chain::resource_limits::resource_limits_object"
	SectionNameResourceUsageObject         SectionName = "eosio::chain::resource_limits::resource_usage_object"
	SectionNameResourceLimitsStateObject   SectionName = "eosio::chain::resource_limits::resource_limits_state_object"
	SectionNameResourceLimitsConfigObject  SectionName = "eosio::chain::resource_limits::resource_limits_config_object"
	SectionNameGenesisState                SectionName = "eosio::chain::genesis_state"
)

func stringToSectionName(name string) SectionName {
	switch name {
	case "eosio::chain::chain_snapshot_header":
		return SectionNameChainSnapshotHeader
	case "eosio::chain::block_state":
		return SectionNameBlockState
	case "eosio::chain::account_object":
		return SectionNameAccountObject
	case "eosio::chain::account_metadata_object":
		return SectionNameAccountMetadataObject
	case "eosio::chain::account_ram_correction_object":
		return SectionNameAccountRamCorrectionObject
	case "eosio::chain::global_property_object":
		return SectionNameGlobalPropertyObject
	case "eosio::chain::protocol_state_object":
		return SectionNameProtocolStateObject
	case "eosio::chain::dynamic_global_property_object":
		return SectionNameDynamicGlobalPropertyObject
	case "eosio::chain::block_summary_object":
		return SectionNameBlockSummaryObject
	case "eosio::chain::transaction_object":
		return SectionNameTransactionObject
	case "eosio::chain::generated_transaction_object":
		return SectionNameGeneratedTransactionObject
	case "eosio::chain::code_object":
		return SectionNameCodeObject
	case "contract_tables":
		return SectionNameContractTables
	case "eosio::chain::permission_object":
		return SectionNamePermissionObject
	case "eosio::chain::permission_link_object":
		return SectionNamePermissionLinkObject
	case "eosio::chain::resource_limits::resource_limits_object":
		return SectionNameResourceLimitsObject
	case "eosio::chain::resource_limits::resource_usage_object":
		return SectionNameResourceUsageObject
	case "eosio::chain::resource_limits::resource_limits_state_object":
		return SectionNameResourceLimitsStateObject
	case "eosio::chain::resource_limits::resource_limits_config_object":
		return SectionNameResourceLimitsConfigObject
	case "eosio::chain::genesis_state":
		return SectionNameGenesisState
	default:
		panic("unsupported section name: " + name)
	}
}

type Section struct {
	Name       SectionName
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
	case SectionNameChainSnapshotHeader:
		return s.readChainSnapshotHeader(f)
	case SectionNameBlockState:
		return s.readBlockState(f)
	case SectionNameAccountObject:
		return s.readAccountObjects(f)
	case SectionNameAccountMetadataObject:
		return s.readAccountMetadataObjects(f)
	case SectionNameAccountRamCorrectionObject:
		return s.readAccountRAMCorrectionObject(f)
	case SectionNameGlobalPropertyObject:
		return s.readGlobalPropertyObject(f)
	case SectionNameProtocolStateObject:
		return s.readProtocolStateObject(f)
	case SectionNameDynamicGlobalPropertyObject:
		return s.readDynamicGlobalPropertyObject(f)
	case SectionNameBlockSummaryObject:
		return s.readBlockSummary(f)
	case SectionNameTransactionObject:
		return s.readTransactionObject(f)
	case SectionNameGeneratedTransactionObject:
		return s.readGeneratedTransactionObject(f)
	case SectionNameCodeObject:
		return s.readCodeObject(f)
	case SectionNameContractTables:
		return s.readContractTables(f)
	case SectionNamePermissionObject:
		return s.readPermissionObject(f)
	case SectionNamePermissionLinkObject:
		return s.readPermissionLinkObject(f)
	case SectionNameResourceLimitsObject:
		return s.readResourceLimitsObject(f)
	case SectionNameResourceUsageObject:
		return s.readResourceUsageObject(f)
	case SectionNameResourceLimitsStateObject:
		return s.readResourceLimitsStateObject(f)
	case SectionNameResourceLimitsConfigObject:
		return s.readResourceLimitsConfigObject(f)
	case SectionNameGenesisState:
		return s.readGenesisState(f)
	default:
		panic("unsupported section: " + s.Name)
	}

}
