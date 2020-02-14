package ship

import (
	"fmt"

	"github.com/eoscanada/eos-go"
)

func NewGetBlocksAck(num uint32) []byte {
	myReq := &Request{
		BaseVariant: eos.BaseVariant{
			TypeID: GetBlocksAckRequestV0Type,
			Impl: &GetBlocksAckRequestV0{
				NumMessages: num,
			},
		},
	}
	bytes, err := eos.MarshalBinary(myReq)
	if err != nil {
		panic(err)
	}

	return bytes
}

func NewRequest(req *GetBlocksRequestV0) []byte {
	myReq := &Request{
		BaseVariant: eos.BaseVariant{
			TypeID: GetBlocksRequestV0Type,
			Impl:   req,
		},
	}
	bytes, err := eos.MarshalBinary(myReq)
	if err != nil {
		panic(err)
	}

	return bytes
}

func ParseGetBlockResultV0(in []byte) (*GetBlocksResultV0, error) {
	variant := &Result{}
	if err := eos.UnmarshalBinary(in, &variant); err != nil {
		return nil, err
	}
	if variant.TypeID != GetBlocksResultV0Type {
		return nil, fmt.Errorf("invalid response type: %d", variant.TypeID)
	}
	return variant.Impl.(*GetBlocksResultV0), nil
}
