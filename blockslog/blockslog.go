package blockslog

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/eoscanada/eos-go"
)

func Process(filename string) error {
	fl, err := os.Open(filename)
	if err != nil {
		return err
	}

	versionData := make([]byte, 4)
	_, err = fl.Read(versionData)
	if err != nil {
		return err
	}

	version := versionData[0]

	fmt.Println("Version", version)

	firstBlockData := []byte{1, 0, 0, 0}

	if version > 1 {
		fmt.Println("Reading first block")
		_, err = fl.Read(firstBlockData)
		if err != nil {
			return err
		}
	}

	firstBlockNum := binary.LittleEndian.Uint32(firstBlockData)
	fmt.Println("First block", firstBlockNum)

	// Certain conditions where the genesis state is written:
	// bool block_log::contains_genesis_state(uint32_t version, uint32_t first_block_num) {
	//    return version <= 2 || first_block_num == 1;
	// }

	chainID := make([]byte, 32)
	if version >= 3 && firstBlockNum > 1 {
		fmt.Println("Reading Chain ID")
		_, err = fl.Read(chainID)
		if err != nil {
			return err
		}
	}

	fmt.Println("Chain ID", hex.EncodeToString(chainID))

	if version != 1 {
		totem := make([]byte, 8)
		_, err := fl.Read(totem)
		if err != nil {
			return err
		}
		fmt.Println("Totem", totem)
	}

	for i := 0 ; i < 5; i++{
		cnt := make([]byte, 1000000)

		prevPos, err := fl.Seek(0, os.SEEK_CUR)
		if err != nil {
			return err
		}

		_, err = fl.Read(cnt)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		d := eos.NewDecoder(cnt)
		var block *eos.SignedBlock

		if err := d.Decode(&block); err != nil {
			return fmt.Errorf("decoding signed block: %w", err)
		}

		jsonStr, err := json.Marshal(block)
		if err != nil {
			return err
		}

		fmt.Println(string(jsonStr))

		if _, err = fl.Seek(prevPos+int64(d.LastPos())+8, os.SEEK_SET); err != nil {
			return err
		}
		fmt.Println("Last pos", d.LastPos())

	}
	return nil
}
