package eosapi

type Transfer struct {
	From     AccountName `json:"from"`
	To       AccountName `json:"to"`
	Quantity uint64      `json:"quantity" struc:"uint64,little"`
	MemoLen  uint32      `json:"-" struc:"varint32,sizeof=Memo"`
	Memo     string      `json:"memo"`
}

type Issue struct {
	To       AccountName `json:"to"`
	Quantity uint64      `json:"quantity" struc:"uint64,little"`
}

type SetCode struct {
	Account   AccountName `json:"account"`
	VMType    byte        `json:"vmtype"`
	VMVersion byte        `json:"vmversion"`
	Code      HexBytes    `json:"bytes"`
}
