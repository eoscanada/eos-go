package system

import "github.com/eoscanada/eos-go"

// NewNonce returns a `nonce` action that lives on the
// `eosio.bios` contract. It should exist only when booting a new
// network, as it is replaced using the `eos-bios` boot process by the
// `eosio.system` contract.
func NewVoteProducer(voter eos.AccountName, proxy eos.AccountName, producers ...eos.AccountName) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("voteproducer"),
		Authorization: []eos.PermissionLevel{
			{Actor: voter, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(
			VoteProducer{
				Voter:     voter,
				Proxy:     proxy,
				Producers: producers,
			},
		),
	}
	return a
}
