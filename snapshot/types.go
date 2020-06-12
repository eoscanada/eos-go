package snapshot

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
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

type KeyValueObject struct {
	PrimKey string
	Payer   string
	Value   []byte
}

func readContractTables(buf *bufio.Reader) error {
	for {
		head := make([]byte, 8+8+8+8+4)
		_, err := buf.Read(head)
		if err != nil {
			return err
		}
		// TODO: check the silenced `written` retval, make sure it equals 36

		t := &TableIDObject{
			Code:      eos.NameToString(binary.LittleEndian.Uint64(head[0:8])),
			Scope:     eos.NameToString(binary.LittleEndian.Uint64(head[8:16])),
			TableName: eos.NameToString(binary.LittleEndian.Uint64(head[16:24])),
			Payer:     eos.NameToString(binary.LittleEndian.Uint64(head[24:32])),
			Count:     binary.LittleEndian.Uint32(head[32:36]),
		}

		size, err := binary.ReadUvarint(buf)
		if err != nil {
			return err
		}

		size32 := uint32(size)

		if t.Count != size32 {
			fmt.Println("WARNING, size and count are NOT equal", t.Count, size)
		}

		fmt.Println("Table:", t.Code, t.Scope, t.TableName, t.Payer, t.Count, size32)

		for i := 0; i < int(size); i++ {
			head := make([]byte, 8+8)
			_, err := buf.Read(head)
			if err != nil {
				return err
			}
			// TODO: check the silenced `written` retval, make sure it equals 16

			kv := &KeyValueObject{
				PrimKey: eos.NameToString(binary.LittleEndian.Uint64(head[0:8])),
				Payer:   eos.NameToString(binary.LittleEndian.Uint64(head[8:16])),
			}

			valueSize, err := binary.ReadUvarint(buf)
			if err != nil {
				return err
			}

			val := make([]byte, valueSize)

			if _, err = buf.Read(val); err != nil {
				// TODO: check the `written`, make sure it equals `valueSize`
				return err
			}

			fmt.Println("  Row:", kv.PrimKey, kv.Payer, len(val))
			kv.Value = val
		}

		skipper := make([]byte, 5)
		if _, err := buf.Read(skipper); err != nil {
			return err
		}

		// TODO: assert it's always 0x0000000000

		fmt.Println("End section", hex.EncodeToString(skipper))

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

type AccountObject struct {
	Name         eos.AccountName
	CreationDate eos.BlockTimestamp
	RawABI       []byte
}

func readAccountObjects(buf *bufio.Reader, count uint64) error {
	for i := uint64(0); i < count; i++ {
		a := AccountObject{}
		cnt := make([]byte, 12)
		_, err := buf.Read(cnt)
		if err != nil {
			return err
		}

		if err := eos.UnmarshalBinary(cnt[:8], &a.Name); err != nil {
			return err
		}

		if err := eos.UnmarshalBinary(cnt[8:12], &a.CreationDate); err != nil {
			return err
		}

		val, err := readByteArray(buf)
		if err != nil {
			return err
		}

		a.RawABI = val

		fmt.Println("Account", a.Name, "created", a.CreationDate.Format(time.RFC3339), "abi length", len(val))
	}
	return nil
}

func readByteArray(buf *bufio.Reader) ([]byte, error) {
	valueSize, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}

	val := make([]byte, valueSize)

	written, err := buf.Read(val)
	if err != nil {
		return nil, err
	}
	if uint64(written) != valueSize {
		return nil, fmt.Errorf("inconsistent read, expected %d bytes, read %d", valueSize, written)
	}

	return val, nil
}
