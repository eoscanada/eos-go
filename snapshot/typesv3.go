package snapshot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/eoscanada/eos-go"
)

type TableIDObject struct {
	Code      string
	Scope     string
	TableName string
	Payer     string
	Count     uint32 // represents the number of rows & indices for a given table
}

type ContractRow struct {
	PrimKey string
	Payer   string
}

type KeyValueObject struct {
	ContractRow
	Value eos.HexBytes
}

type Index64Object struct {
	ContractRow
	SecondaryKey eos.Name
}

type Index128Object struct {
	ContractRow
	SecondaryKey eos.Uint128
}

type Index256Object struct {
	ContractRow
	SecondaryKey eos.Checksum256
}

type IndexDoubleObject struct {
	ContractRow
	SecondaryKey eos.Float64
}

type IndexLongDoubleObject struct {
	ContractRow
	SecondaryKey eos.Float128
}

func (section *Section) readContractTables(f callbackFunc) error {
	fl := section.Buffer

	bufSize := section.BufferSize
	bytesBuf := make([]byte, bufSize)
	slurped, err := fl.Read(bytesBuf)
	if err != nil {
		return err
	}
	if slurped != int(bufSize) {
		slurped2, err := fl.Read(bytesBuf[slurped:])
		if err != nil {
			return err
		}
		if slurped+slurped2 != int(bufSize) {
			return fmt.Errorf("read less than section size: %d of %d", slurped+slurped2, section.BufferSize)
		}
	}
	buf := bytes.NewBuffer(bytesBuf)

	for {
		head := make([]byte, 8+8+8+8+4)
		readz, err := buf.Read(head)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading table id: %w", err)
		}
		if readz != 8+8+8+8+4 {
			return fmt.Errorf("incomplete read for table_id_object: %d out of %d", readz, 8+8+8+8+4)
		}

		t := &TableIDObject{
			Code:      eos.NameToString(binary.LittleEndian.Uint64(head[0:8])),
			Scope:     eos.NameToString(binary.LittleEndian.Uint64(head[8:16])),
			TableName: eos.NameToString(binary.LittleEndian.Uint64(head[16:24])),
			Payer:     eos.NameToString(binary.LittleEndian.Uint64(head[24:32])),
			Count:     binary.LittleEndian.Uint32(head[32:36]),
		}

		if err := f(t); err != nil {
			return err
		}

		// if t.Code == "houseaccount" {
		// 	fmt.Println("Table code", t.Code, "scope", t.Scope, "tablename", t.TableName, "payer", t.Payer, "count", t.Count)
		// }

		for idxType := 0; idxType < 6; idxType++ {

			// offset, _ := buf.(*os.File).Seek(0, os.SEEK_CUR)
			// fmt.Println("OFFSET", offset)

			size, err := readUvarint(buf)
			if err != nil {
				return fmt.Errorf("reading index type size: %w", err)
			}

			// offset2, _ := buf.(*os.File).Seek(0, os.SEEK_CUR)
			// fmt.Println("OFFSET n size", offset2, offset2-offset, size)

			// if t.Code == "houseaccount" {
			// 	fmt.Println("  index type", idxType, "index size", size, "code", t.Code, "scope", t.Scope, "tablename", t.TableName, "count", t.Count)
			// }

			for i := 0; i < int(size); i++ {
				head := make([]byte, 8+8)
				readz, err := buf.Read(head)
				if err != nil {
					return fmt.Errorf("reading key value head: %w", err)
				}
				if readz != 16 {
					return fmt.Errorf("incomplete read for row header: %d out of 16", readz)
				}

				contractRow := ContractRow{
					PrimKey: eos.NameToString(binary.LittleEndian.Uint64(head[0:8])),
					Payer:   eos.NameToString(binary.LittleEndian.Uint64(head[8:16])),
				}

				var row interface{}
				switch idxType {

				case 0: /* key_value_object */
					obj := &KeyValueObject{ContractRow: contractRow}

					valueSize, err := readUvarint(buf)
					if err != nil {
						return err
					}
					val := make([]byte, valueSize)
					readz, err = buf.Read(val)
					if err != nil {
						return err
					}
					if readz != int(valueSize) {
						return fmt.Errorf("incomplete read key_value_object: %d out of %d", readz, valueSize)
					}

					obj.Value = val
					row = obj

				case 1: /* index64_object */
					obj := &Index64Object{ContractRow: contractRow}
					val := make([]byte, 8)
					readz, err := buf.Read(val)
					if err != nil {
						return err
					}
					if readz != 8 {
						return fmt.Errorf("incomplete read index64_object: %d out of 8", readz)
					}
					if err := eos.UnmarshalBinary(val, &obj.SecondaryKey); err != nil {
						return err
					}
					row = obj

				case 2: /* index128_object */
					obj := &Index128Object{ContractRow: contractRow}
					val := make([]byte, 16)
					if _, err = buf.Read(val); err != nil {
						return err
					}
					if err := eos.UnmarshalBinary(val, &obj.SecondaryKey); err != nil {
						return err
					}
					row = obj
				case 3: /* index256_object */
					obj := &Index256Object{ContractRow: contractRow}
					val := make([]byte, 32)
					if _, err = buf.Read(val); err != nil {
						return err
					}
					if err := eos.UnmarshalBinary(val, &obj.SecondaryKey); err != nil {
						return err
					}
					row = obj
				case 4: /* index_double_object */
					obj := &IndexDoubleObject{ContractRow: contractRow}
					val := make([]byte, 8)
					if _, err = buf.Read(val); err != nil {
						return err
					}
					if err := eos.UnmarshalBinary(val, &obj.SecondaryKey); err != nil {
						return err
					}
					row = obj
				case 5: /* index_long_double_object */
					obj := &IndexLongDoubleObject{ContractRow: contractRow}
					val := make([]byte, 16)
					if _, err = buf.Read(val); err != nil {
						return err
					}
					if err := eos.UnmarshalBinary(val, &obj.SecondaryKey); err != nil {
						return err
					}
					row = obj
				}

				// if t.Code == "houseaccount" {
				// 	out, _ := json.Marshal(row)
				// 	fmt.Printf("%T: %d. %s seek: %d\n", row, i, string(out), section.BufferSize)
				// 	//fmt.Printf("Buffer state: size %d buffered %d\n", section.Buffer.Size())
				// }

				if err := f(row); err != nil {
					return err
				}

			}

		}
	}

	return nil
}

// type TableIDObject struct {
//    class table_id_object : public chainbase::object<table_id_object_type, table_id_object> {
//       OBJECT_CTOR(table_id_object)

//       id_type        id;
//       account_name   code;  //< code should not be changed within a chainbase modifier lambda
//       scope_name     scope; //< scope should not be changed within a chainbase modifier lambda
//       table_name     table; //< table should not be changed within a chainbase modifier lambda
//       account_name   payer;
//       uint32_t       count = 0; /// the number of elements in the table
//    };
// }
// type KeyValueObject struct {
// 	   struct key_value_object : public chainbase::object<key_value_object_type, key_value_object> {
//       OBJECT_CTOR(key_value_object, (value))

//       typedef uint64_t key_type;
//       static const int number_of_keys = 1;

//       id_type               id;
//       table_id              t_id; //< t_id should not be changed within a chainbase modifier lambda
//       uint64_t              primary_key; //< primary_key should not be changed within a chainbase modifier lambda
//       account_name          payer;
//       shared_blob           value;
//    };
// }

type BlockState struct {
	/// from block_header_state_common
	BlockNum                         uint32                         `json:"block_num"`
	DposProposedIrreversibleBlocknum uint32                         `json:"dpos_proposed_irreversible_blocknum"`
	DposIrreversibleBlocknum         uint32                         `json:"dpos_irreversible_blocknum"`
	ActiveSchedule                   *eos.ProducerAuthoritySchedule `json:"active_schedule"`
	BlockrootMerkle                  *eos.MerkleRoot                `json:"blockroot_merkle"`
	ProducerToLastProduced           []eos.PairAccountNameBlockNum  `json:"producer_to_last_produced"`
	ProducerToLastImpliedIrb         []eos.PairAccountNameBlockNum  `json:"producer_to_last_implied_irb"`
	BlockSigningKey                  *eos.BlockSigningAuthority     `json:"block_signing_key"`
	ConfirmCount                     []uint8                        `json:"confirm_count"`

	// from block_header_state
	BlockID                   eos.Checksum256                   `json:"id"`
	Header                    *eos.SignedBlockHeader            `json:"header"`
	PendingSchedule           *ScheduleInfo                     `json:"pending_schedule"`
	ActivatedProtocolFeatures *eos.ProtocolFeatureActivationSet `json:"activated_protocol_features"`
}

type ScheduleInfo struct {
	ScheduleLIBNum uint32                         `json:"schedule_lib_num"`
	ScheduleHash   eos.Checksum256                `json:"schedule_hash"`
	Schedule       *eos.ProducerAuthoritySchedule `json:"schedule"`
}

func (section *Section) readBlockState(f callbackFunc) (err error) {
	cnt := make([]byte, section.BufferSize)
	_, err = section.Buffer.Read(cnt)
	if err != nil {
		return
	}

	var state BlockState
	err = eos.UnmarshalBinary(cnt, &state)
	if err != nil {
		return
	}

	if err := f(state); err != nil {
		return err
	}

	return nil
}

////

type AccountObject struct {
	Name         eos.AccountName
	CreationDate eos.BlockTimestamp
	RawABI       []byte
}

func (section *Section) readAccountObjects(f callbackFunc) error {
	for i := uint64(0); i < section.RowCount; i++ {
		a := AccountObject{}
		cnt := make([]byte, 12)
		_, err := section.Buffer.Read(cnt)
		if err != nil {
			return err
		}

		if err := eos.UnmarshalBinary(cnt[:8], &a.Name); err != nil {
			return err
		}

		if err := eos.UnmarshalBinary(cnt[8:12], &a.CreationDate); err != nil {
			return err
		}

		val, err := readByteArray(section.Buffer)
		if err != nil {
			return err
		}

		a.RawABI = val

		if err := f(a); err != nil {
			return err
		}
		//fmt.Println("Account", a.Name, "created", a.CreationDate.Format(time.RFC3339), "abi length", len(val))
	}
	return nil
}

////

type AccountMetadataObject struct {
	Name           eos.AccountName //< name should not be changed within a chainbase modifier lambda
	RecvSequence   eos.Uint64
	AuthSequence   eos.Uint64
	CodeSequence   eos.Uint64
	ABISequence    eos.Uint64
	CodeHash       eos.Checksum256
	LastCodeUpdate eos.TimePoint
	Flags          uint32 // First flag means "privileged".
	VMType         byte
	VMVersion      byte
}

func (section *Section) readAccountMetadataObjects(f callbackFunc) error {
	for i := uint64(0); i < section.RowCount; i++ {
		a := AccountMetadataObject{}
		cnt := make([]byte, 86) // account_metadata_object is fixed size 86 bytes
		_, err := section.Buffer.Read(cnt)
		if err != nil {
			return err
		}

		if err := eos.UnmarshalBinary(cnt, &a); err != nil {
			return err
		}

		if err := f(a); err != nil {
			return err
		}
		// fmt.Println("Account", a.Name, "ast code updatecreated", a.LastCodeUpdate, a.RecvSequence, a.AuthSequence, a.CodeSequence, a.ABISequence, a.VMType, a.VMVersion, "flags", a.Flags)
	}
	return nil
}

////

type ChainSnapshotHeader struct {
	Version uint32
}

func (section *Section) readChainSnapshotHeader(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	_, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}

	var header ChainSnapshotHeader
	err = eos.UnmarshalBinary(cnt, &header)
	if err != nil {
		return err
	}

	if err := f(header); err != nil {
		return err
	}

	return nil
}

type GlobalPropertyObject struct {
	ProposedScheduleBlockNum uint32 `eos:"optional"`
	ProposedSchedule         *eos.ProducerAuthoritySchedule
	Configuration            ChainConfig
	ChainID                  eos.Checksum256
}

func (section *Section) readGlobalPropertyObject(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	_, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}

	var obj GlobalPropertyObject
	err = eos.UnmarshalBinary(cnt, &obj)
	if err != nil {
		return err
	}

	if err := f(obj); err != nil {
		return err
	}

	return nil
}

//

type ProtocolStateObject struct {
	ActivatedProtocolFeatures    []*ActivatedProtocolFeature
	PreactivatedProtocolFeatures []eos.Checksum256
	WhitelistedIntrinsics        []string
	NumSupportedKeyTypes         uint32
}

type ActivatedProtocolFeature struct {
	FeatureDigest      eos.Checksum256
	ActivationBlockNum uint32
}

func (section *Section) readProtocolStateObject(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	_, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}

	// _ = ioutil.WriteFile("/tmp/test.dat", cnt, 0664)

	var obj ProtocolStateObject
	err = eos.UnmarshalBinary(cnt, &obj)
	if err != nil {
		return err
	}

	if err := f(obj); err != nil {
		return err
	}

	return nil
}

//

type DynamicGlobalPropertyObject struct {
	GlobalActionSequence eos.Uint64
}

func (section *Section) readDynamicGlobalPropertyObject(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	_, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}

	var obj DynamicGlobalPropertyObject
	err = eos.UnmarshalBinary(cnt, &obj)
	if err != nil {
		return err
	}

	if err := f(obj); err != nil {
		return err
	}

	return nil
}

//

type AccountRAMCorrectionObject struct {
	Name          eos.AccountName
	RAMCorrection eos.Uint64
}

func (section *Section) readAccountRAMCorrectionObject(f callbackFunc) error {
	for i := uint64(0); i < section.RowCount; i++ {
		a := AccountRAMCorrectionObject{}
		cnt := make([]byte, 16) // fixed size of account_ram_correction_object
		_, err := section.Buffer.Read(cnt)
		if err != nil {
			return err
		}

		if err := eos.UnmarshalBinary(cnt, &a); err != nil {
			return err
		}

		if err := f(a); err != nil {
			return err
		}
	}
	return nil
}

//

type BlockSummary struct {
	BlockID eos.Checksum256
}

func (section *Section) readBlockSummary(f callbackFunc) error {
	for i := uint64(0); i < section.RowCount; i++ {
		a := BlockSummary{}
		cnt := make([]byte, 32) // fixed size of block_summary
		_, err := section.Buffer.Read(cnt)
		if err != nil {
			return err
		}

		if err := eos.UnmarshalBinary(cnt, &a); err != nil {
			return err
		}

		if err := f(a); err != nil {
			return err
		}
	}
	return nil
}

///

type PermissionObject struct { /* special snapshot version of the object */
	Parent      eos.PermissionName ///< parent permission
	Owner       eos.AccountName    ///< the account this permission belongs to
	Name        eos.PermissionName ///< human-readable name for the permission
	LastUpdated eos.TimePoint      ///< the last time this authority was updated
	LastUsed    eos.TimePoint      ///< when this permission was last used
	Auth        eos.Authority      ///< authority required to execute this permission
}

func (section *Section) readPermissionObject(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	_, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}

	for pos := 0; pos < int(section.BufferSize); {
		d := eos.NewDecoder(cnt[pos:])
		var po PermissionObject
		err = d.Decode(&po)
		if err != nil {
			return err
		}

		if err := f(po); err != nil {
			return err
		}

		pos += d.LastPos()
	}

	// _ = ioutil.WriteFile("/tmp/test.dat", cnt, 0664)

	return nil
}

///

type PermissionLinkObject struct {
	/// The account which is defining its permission requirements
	Account eos.AccountName
	/// The contract which account requires @ref required_permission to invoke
	Code eos.AccountName
	/// The message type which account requires @ref required_permission to invoke
	/// May be empty; if so, it sets a default @ref required_permission for all messages to @ref code
	MessageType eos.ActionName
	/// The permission level which @ref account requires for the specified message types
	/// all of the above fields should not be changed within a chainbase modifier lambda
	RequiredPermission eos.PermissionName
}

func (section *Section) readPermissionLinkObject(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	_, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}

	for pos := 0; pos < int(section.BufferSize); {
		d := eos.NewDecoder(cnt[pos:])
		var po PermissionLinkObject
		err = d.Decode(&po)
		if err != nil {
			return err
		}

		if err := f(po); err != nil {
			return err
		}

		pos += d.LastPos()
	}

	// _ = ioutil.WriteFile("/tmp/test.dat", cnt, 0664)

	return nil
}

////

type ResourceLimitsObject struct {
	Owner eos.AccountName //<  should not be changed within a chainbase modifier lambda

	NetWeight eos.Int64
	CPUWeight eos.Int64
	RAMBytes  eos.Int64
}

func (section *Section) readResourceLimitsObject(f callbackFunc) error {
	for i := uint64(0); i < section.RowCount; i++ {
		a := ResourceLimitsObject{}
		cnt := make([]byte, 8+8+8+8) // fixed size of resource_limits_object
		_, err := section.Buffer.Read(cnt)
		if err != nil {
			return err
		}

		if err := eos.UnmarshalBinary(cnt, &a); err != nil {
			return err
		}

		if err := f(a); err != nil {
			return err
		}
	}
	return nil
}

////

type ResourceUsageObject struct {
	Owner eos.AccountName //< owner should not be changed within a chainbase modifier lambda

	NetUsage UsageAccumulator
	CPUUsage UsageAccumulator

	RAMUsage eos.Uint64
}

func (section *Section) readResourceUsageObject(f callbackFunc) error {
	for i := uint64(0); i < section.RowCount; i++ {
		a := ResourceUsageObject{}
		cnt := make([]byte, 8+20+20+8) // fixed size of resource_limits_object
		_, err := section.Buffer.Read(cnt)
		if err != nil {
			return err
		}

		if err := eos.UnmarshalBinary(cnt, &a); err != nil {
			return err
		}

		if err := f(a); err != nil {
			return err
		}
	}
	return nil
}

////

type ResourceLimitsStateObject struct {
	/**
	 * Track the average netusage for blocks
	 */
	AverageBlockNetUsage UsageAccumulator

	/**
	 * Track the average cpu usage for blocks
	 */
	AverageBlockCPUUsage UsageAccumulator

	PendingNetUsage eos.Uint64
	PendingCPUUsage eos.Uint64

	TotalNetWeight eos.Uint64
	TotalCPUWeight eos.Uint64
	TotalRAMBytes  eos.Uint64

	/**
	 * The virtual number of bytes that would be consumed over blocksize_average_window_ms
	 * if all blocks were at their maximum virtual size. This is virtual because the
	 * real maximum block is less, this virtual number is only used for rate limiting users.
	 *
	 * It's lowest possible value is max_block_size * blocksize_average_window_ms / block_interval
	 * It's highest possible value is config::maximum_elastic_resource_multiplier (1000) times its lowest possible value
	 *
	 * This means that the most an account can consume during idle periods is 1000x the bandwidth
	 * it is gauranteed under congestion.
	 *
	 * Increases when average_block_size < target_block_size, decreases when
	 * average_block_size > target_block_size, with a cap at 1000x max_block_size
	 * and a floor at max_block_size;
	 **/
	VirtualNetLimit eos.Uint64

	/**
	 *  Increases when average_bloc
	 */
	VirtualCPULimit eos.Uint64
}

type UsageAccumulator struct {
	LastOrdinal uint32     ///< The ordinal of the last period which has contributed to the average
	ValueEx     eos.Uint64 ///< The current average pre-multiplied by Precision
	Consumed    eos.Uint64 ///< The last periods average + the current periods contribution so far
}

func (section *Section) readResourceLimitsStateObject(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	_, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}

	// _ = ioutil.WriteFile("/tmp/test.dat", cnt, 0664)

	var obj ResourceLimitsStateObject
	err = eos.UnmarshalBinary(cnt, &obj)
	if err != nil {
		return err
	}

	if err := f(obj); err != nil {
		return err
	}

	return nil
}

////

type ResourceLimitsConfigObject struct {
	CPULimitParameters ElasticLimitParameters
	NetLimitParameters ElasticLimitParameters

	AccountCPUUsageAverageWindow uint32
	AccountNetUsageAverageWindow uint32
}

type ElasticLimitParameters struct {
	Target  eos.Uint64 // the desired usage
	Max     eos.Uint64 // the maximum usage
	Periods uint32     // the number of aggregation periods that contribute to the average usage

	MaxMultiplier uint32     // the multiplier by which virtual space can oversell usage when uncongested
	ContractRate  eos.Uint64 // the rate at which a congested resource contracts its limit
	ExpandRate    eos.Uint64 // the rate at which an uncongested resource expands its limits
}

func (section *Section) readResourceLimitsConfigObject(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	_, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}

	// _ = ioutil.WriteFile("/tmp/test.dat", cnt, 0664)

	var obj ResourceLimitsConfigObject
	err = eos.UnmarshalBinary(cnt, &obj)
	if err != nil {
		return err
	}

	if err := f(obj); err != nil {
		return err
	}

	return nil
}

////

type CodeObject struct {
	CodeHash       eos.Checksum256 //< code_hash should not be changed within a chainbase modifier lambda
	Code           eos.HexBytes
	CodeRefCount   eos.Uint64
	FirstBlockUsed uint32
	VMType         uint8 //< vm_type should not be changed within a chainbase modifier lambda
	VMVersion      uint8 //< vm_version should not be changed within a chainbase modifier lambda
}

func (section *Section) readCodeObject(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	readz, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}
	if readz != len(cnt) {
		return fmt.Errorf("failed reading the whole code object section: %d of %d", readz, len(cnt))
	}

	for pos := 0; pos < int(section.BufferSize); {
		d := eos.NewDecoder(cnt[pos:])
		var co CodeObject
		err = d.Decode(&co)
		if err != nil {
			return err
		}

		if err := f(co); err != nil {
			return err
		}

		pos += d.LastPos()
	}

	return nil
}

////

type GeneratedTransactionObject struct {
	TrxID      eos.Checksum256 //< trx_id should not be changed within a chainbase modifier lambda
	Sender     eos.AccountName //< sender should not be changed within a chainbase modifier lambda
	SenderID   eos.Uint128     /// ID given this transaction by the sender (should not be changed within a chainbase modifier lambda)
	Payer      eos.AccountName
	DelayUntil eos.TimePoint /// this generated transaction will not be applied until the specified time
	Expiration eos.TimePoint /// this generated transaction will not be applied after  time
	Published  eos.TimePoint
	PackedTrx  eos.HexBytes
}

func (section *Section) readGeneratedTransactionObject(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	readz, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}
	if readz != len(cnt) {
		return fmt.Errorf("failed reading the whole code object section: %d of %d", readz, len(cnt))
	}

	for pos := 0; pos < int(section.BufferSize); {
		d := eos.NewDecoder(cnt[pos:])
		var gto GeneratedTransactionObject
		err = d.Decode(&gto)
		if err != nil {
			return err
		}

		if err := f(gto); err != nil {
			return err
		}

		pos += d.LastPos()
	}

	return nil
}

////

type TransactionObject struct {
	Expiration eos.TimePointSec
	TrxID      eos.Checksum256 //< trx_id shou
}

func (section *Section) readTransactionObject(f callbackFunc) error {
	cnt := make([]byte, section.BufferSize)
	readz, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}
	if readz != len(cnt) {
		return fmt.Errorf("failed reading the whole code object section: %d of %d", readz, len(cnt))
	}

	for pos := 0; pos < int(section.BufferSize); {
		d := eos.NewDecoder(cnt[pos:])
		var to TransactionObject
		err = d.Decode(&to)
		if err != nil {
			return err
		}

		if err := f(to); err != nil {
			return err
		}

		pos += d.LastPos()
	}

	return nil
}
