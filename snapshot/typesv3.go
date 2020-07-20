package snapshot

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/eoscanada/eos-go"
)

type TableIDObject struct {
	Code      string
	Scope     string
	TableName string
	Payer     string
	Count     uint32
}

type ContractRow struct {
	PrimKey string
	Payer   string
}

type KeyValueObject struct {
	ContractRow
	Value []byte
}

type Index64Object struct {
	ContractRow
	SecondaryKey eos.Uint64
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

func readContractTables(section *Section) error {
	buf := section.Buffer
	for {
		head := make([]byte, 8+8+8+8+4)
		_, err := buf.Read(head)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading table id: %w", err)
		}
		// TODO: check the silenced `written` retval, make sure it equals 36

		t := &TableIDObject{
			Code:      eos.NameToString(binary.LittleEndian.Uint64(head[0:8])),
			Scope:     eos.NameToString(binary.LittleEndian.Uint64(head[8:16])),
			TableName: eos.NameToString(binary.LittleEndian.Uint64(head[16:24])),
			Payer:     eos.NameToString(binary.LittleEndian.Uint64(head[24:32])),
			Count:     binary.LittleEndian.Uint32(head[32:36]),
		}

		fmt.Println("Table:", t.Code, t.Scope, t.TableName, t.Payer, t.Count)

		for idxType := 0; idxType < 6; idxType++ {

			size, err := binary.ReadUvarint(buf)
			if err != nil {
				return fmt.Errorf("reading index type size: %w", err)
			}

			size32 := uint32(size)

			if t.Count != 1 {
				fmt.Println("  index type", idxType, "size", size32, "code", t.Code, "scope", t.Scope, "tablename", t.TableName, "count", t.Count)
			}

			// if t.Count != size32 {
			// 	return fmt.Errorf("WARNING, size and count are NOT equal: %d %d", t.Count, size)
			// }

			for i := 0; i < int(size); i++ {
				head := make([]byte, 8+8)
				_, err := buf.Read(head)
				// TODO: check the silenced `written` retval, make sure it equals 16
				if err != nil {
					return fmt.Errorf("reading key value head: %w", err)
				}

				contractRow := ContractRow{
					PrimKey: eos.NameToString(binary.LittleEndian.Uint64(head[0:8])),
					Payer:   eos.NameToString(binary.LittleEndian.Uint64(head[8:16])),
				}

				var row interface{}
				switch idxType {

				case 0: /* key_value_object */
					obj := &KeyValueObject{ContractRow: contractRow}

					valueSize, err := binary.ReadUvarint(buf)
					if err != nil {
						return err
					}

					val := make([]byte, valueSize)
					if _, err = buf.Read(val); err != nil {
						// TODO: check the `written`, make sure it equals `valueSize`
						return err
					}

					obj.Value = val
					row = obj

				case 1: /* index64_object */
					obj := &Index64Object{ContractRow: contractRow}
					val := make([]byte, 8)
					if _, err = buf.Read(val); err != nil {
						return err
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

				if t.Code == "houseaccount" {
					out, _ := json.Marshal(row)
					fmt.Printf("%T: %s\n", row, string(out))
				}
			}

		}

		// TODO: assert it's always 0x0000000000

		// if bytes.Compare(skipper, []byte{0, 0, 0, 0, 0}) != 0 {
		// 	return fmt.Errorf("failed skipper check, perhaps that means something? %v", skipper)
		// }
		//fmt.Println("End section", hex.EncodeToString(skipper))

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

func readBlockState(section *Section) (err error) {
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
	out, _ := json.MarshalIndent(state, "  ", "  ")
	fmt.Println(string(out))

	// FIXME: could handle the failed `bool` decoding,
	// something BlockHeaderState
	// FC_REFLECT_DERIVED( eosio::chain::block_state, (eosio::chain::block_header_state), (block)(validated) )
	// and the `BlockState` we have in `eos-go`.

	return nil
}

////

type AccountObject struct {
	Name         eos.AccountName
	CreationDate eos.BlockTimestamp
	RawABI       []byte
}

func readAccountObjects(section *Section) error {
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

		fmt.Println("Account", a.Name, "created", a.CreationDate.Format(time.RFC3339), "abi length", len(val))
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

func readAccountMetadataObjects(section *Section) error {
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

		fmt.Println("Account", a.Name, "ast code updatecreated", a.LastCodeUpdate, a.RecvSequence, a.AuthSequence, a.CodeSequence, a.ABISequence, a.VMType, a.VMVersion, "flags", a.Flags)
	}
	return nil
}

////

type ChainSnapshotHeader struct {
	Version uint32
}

func readChainSnapshotHeader(section *Section) error {
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

	cnt, _ = json.MarshalIndent(header, "  ", "  ")
	fmt.Println(string(cnt))

	return nil
}

type GlobalPropertyObject struct {
	ProposedScheduleBlockNum uint32 `eos:"optional"`
	ProposedSchedule         *eos.ProducerAuthoritySchedule
	Configuration            ChainConfig
	ChainID                  eos.Checksum256
}

func readGlobalPropertyObject(section *Section) error {
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

	cnt, _ = json.MarshalIndent(obj, "  ", "  ")
	fmt.Println(string(cnt))

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

func readProtocolStateObject(section *Section) error {
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

	cnt, _ = json.MarshalIndent(obj, "  ", "  ")
	fmt.Println(string(cnt))

	return nil
}

//

type DynamicGlobalPropertyObject struct {
	GlobalActionSequence eos.Uint64
}

func readDynamicGlobalPropertyObject(section *Section) error {
	cnt := make([]byte, section.BufferSize)
	_, err := section.Buffer.Read(cnt)
	if err != nil {
		return err
	}

	// _ = ioutil.WriteFile("/tmp/test.dat", cnt, 0664)

	var obj DynamicGlobalPropertyObject
	err = eos.UnmarshalBinary(cnt, &obj)
	if err != nil {
		return err
	}

	cnt, _ = json.MarshalIndent(obj, "  ", "  ")
	fmt.Println(string(cnt))

	return nil
}

//

type AccountRAMCorrectionObject struct {
	Name          eos.AccountName
	RAMCorrection eos.Uint64
}

func readAccountRAMCorrectionObject(section *Section) error {
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

		fmt.Println("Account", a.Name, a.RAMCorrection)
	}
	return nil
}

//

type BlockSummary struct {
	BlockID eos.Checksum256
}

func readBlockSummary(section *Section) error {
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

		fmt.Println("Block", a.BlockID)
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

func readPermissionObject(section *Section) error {
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

		out, _ := json.MarshalIndent(po, "  ", "  ")
		fmt.Println(string(out))
		// fmt.Println("NEXT", d.LastPos())

		pos += d.LastPos()
	}

	// _ = ioutil.WriteFile("/tmp/test.dat", cnt, 0664)

	return nil
}
