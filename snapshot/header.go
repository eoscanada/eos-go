package snapshot

// Header is the top-most header, which determines the file format.
//
// It is not to be confused with the
// eosio::chain::chain_snapshot_header which talks about the version
// of the contents of the snapshot file.
type Header struct {
	Version uint32
}
