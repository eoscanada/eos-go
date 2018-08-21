package system

import (
	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

// NewRegProducer returns a `regproducer` action that lives on the
// `eosio.system` contract.
func NewRegProducer(producer eos.AccountName, producerKey ecc.PublicKey, url string, location uint16) *eos.Action {
	return &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("regproducer"),
		Authorization: []eos.PermissionLevel{
			{Actor: producer, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(RegProducer{
			Producer:    producer,
			ProducerKey: producerKey,
			URL:         url,
			Location:    location,
		}),
	}
}

// RegProducer represents the `eosio.system::regproducer` action
type RegProducer struct {
	Producer    eos.AccountName `json:"producer"`
	ProducerKey ecc.PublicKey   `json:"producer_key"`
	URL         string          `json:"url"`
	Location    uint16          `json:"location"` // what,s the meaning of that anyway ?
}
