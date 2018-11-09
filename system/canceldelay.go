package system

import "github.com/eoscanada/eos-go"

// NewCancelDelay creates an action from the `eosio.system` contract
// called `canceldelay`.
//
// `canceldelay` allows you to cancel a deferred transaction,
// previously sent to the chain with a `delay_sec` larger than 0.  You
// need to sign with cancelingAuth, to cancel a transaction signed
// with that same authority.
func NewCancelDelay(cancelingAuth eos.PermissionLevel, transactionID eos.Checksum256) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("canceldelay"),
		Authorization: []eos.PermissionLevel{
			cancelingAuth,
		},
		ActionData: eos.NewActionData(CancelDelay{
			CancelingAuth: cancelingAuth,
			TransactionID: transactionID,
		}),
	}

	return a
}

// CancelDelay represents the native `canceldelay` action, through the
// system contract.
type CancelDelay struct {
	CancelingAuth eos.PermissionLevel `json:"canceling_auth"`
	TransactionID eos.Checksum256     `json:"trx_id"`
}
