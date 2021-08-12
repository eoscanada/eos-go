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

	// Ultra Specific
	SectionAccountFreeActionsObject SectionName = "eosio::chain::account_free_actions_object"
)

type Section struct {
	Name       SectionName
	Offset     uint64
	Size       uint64 // This includes the section name and row count
	BufferSize uint64 // This represents the bytes that are following the section header
	RowCount   uint64 // This is a count of rows packed in `Buffer`
	Buffer     io.Reader
}

type sectionHandlerFunc func(s *Section, f sectionCallbackFunc) error
type sectionCallbackFunc func(obj interface{}) error
