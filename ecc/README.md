EOSIO Elliptic Curve Cryptography Wrapper
=========================================

This is a simple wrapper for `btcec`, that handles the specificities
of the format for keys in EOS.

It was crafted in reference to `eosjs-ecc`, `eosjs-keygen` and the C++
codebase of EOS.IO Software.

This handles the `EOS` prefix on public keys, manages the version and
checksums in public and private keys.
