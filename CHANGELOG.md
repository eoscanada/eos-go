## Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

### Unreleased

* Added a check for int overflows in `ReadByteArray`
* Changed valueToInt, valueToUint, valueToFload function in abiencode.go for compatible with double quoted string to number.
* Changed `NewAssetFromString` validation to allow parsing of empty assets
* Added `action_trace_v1` field
* Added `AsTime` helper functions to convert `TimePoint` and `TimePointSec` to `time.Time`
* Added support for decoding action results
* Raised the minimum supporting Go version 1.13 -> 1.17
* Added unit tests of REST API
* Added `earliest_available_block_num`, `fork_db_head_block_num`, `fork_db_head_block_id`, `last_irreversible_block_time`, `total_cpu_weight`, `total_net_weight`, & `server_full_version_string` fields in `InfoResp`
* Added `head_block_num`, `head_block_time`, `rex_info`, `subjective_cpu_bill_limit`, & `eosio_any_linked_actions` fields in `AccountResp`
* Added the new type `RexInfo`, `LinkedAction`
* Added `linked_actions` in `Permission`
* Added `context_free_data` & `transaction` fields in `PackedTransaction`

#### Breaking Changes

* `AccountResp.last_code_update` & `AccountResp.created` in `AccountResp` are now `BlockTimestamp`, were previously `JSONTime`

#### Added

#### Changed

#### Fixed

* Improved the error handling when decoding table rows with variant types.

* Fixed decoding of table rows with variant types.

* Fixed serialization of `map[K]V` when using `eos.MarshalBinary` so that ordering in which the keys are serialized is in lexicographical order, just like JSON serialization.

* Updated to latest version of `github.com/streamingfast/logging` library.

* Fix NewPrivateKey correctly working with PVT_K1 keys (#177)

* Fixed ABI encoder for `SymbolCode` type.

#### Deprecated

### [**0.10.2**](https://github.com/eoscanada/eos-go/releases/tag/v0.10.2) (January 19th, 2022)

#### Changed

* Changed ABI encoder `varuint32` to re-use existing encoder.

#### Fixed

* Fixed ABI encoder `varint32` types.

### [**0.10.1**](https://github.com/eoscanada/eos-go/releases/tag/v0.10.1) (January 19th, 2022)

#### Breaking Changes

* Renamed [`BlockTimestampFormat`](https://github.com/eoscanada/eos-go/blob/a1623cc5a2223005a4dc7d4dec972d6119de42ff/types.go#L844) to `blockTimestampFormat` making it private.

#### Added

* Add EOSIO Testnet symbol (#165)

* Update ActionResp (#166)

#### Changed

* Updated to latest version of `github.com/streamingfast/logging` library.

* `eos.MarshalBinary` will now refuses to serialize a `map[K]V` if `K`'s type is not comparable.

#### Fixed

* Fixed serialization of `map[K]V` types by sorting the keys to have a stable order.

* Bugfix StringToSymbol (44b6fbd)

* Fixed built-in examples (pointing by default to EOS Nation API nodes) (36114bd)

#### Deprecated

### [**0.10**](https://github.com/eoscanada/eos-go/releases/tag/v0.10.0) (October 16th, 2021)

#### Breaking Changes
- **BREAKING**: We started adding an initial `context.Context` to all interruptible functions. All method performing an HTTP call have the new parameter as well as a bunch of other method. We cannot list all of them. If the caller already have a `context.Context` value, pass it to the function that now require one. Otherwise, simply pass `context.Background()`.

- **BREAKING**: Fixed binary unpacking of `BlockState`, `TransactionTrace`, `SignedTransaction`, `Action` (and some inner types). This required changing a few struct fields to better fit with EOSIO definition, here the full list:
  - `MerkleRoot.ActiveNodes` is now a `[]Checksum256`, was previously `[]string`
  - `MerkleRoot.NodeCount` is now a `uint64`, was previously `uint32`
  - Type `EOSNameOrUint32` has been removed and replaced by `PairAccountNameBlockNum` which is strictly typed now.
  - `BlockState.ProducerToLastProduced` is now `[]PairAccountNameBlockNum`, was previously `[][2]EOSNameOrUint32`.
  - `BlockState.ProducerToLastImpliedIRB` is now `[]PairAccountNameBlockNum`, was previously `[][2]EOSNameOrUint32`.
  - `BlockState.BlockID` is now a `Checksum256`, was previously `string`.
  - `BlockState.ActivatedProtocolFeatures` is now a `*ProtocolFeatureActivationSet`, was previously `map[string][]HexBytes`.
  - `BlockState.ConfirmCount` is now a `[]uint8`, was previously `[]uint32`.
  - `PendingSchedule.ScheduleHash` is now a `Checksum256`, was previously `HexBytes`.
  - `ActionTraceReceipt.ActionDigest` is now a `Checksum256`, was previously `string`.
  - `ActionTraceReceipt.CodeSequence` is now a `Varuint32`, was previously `Uint64`.
  - `ActionTraceReceipt.ABISequence` is now a `Varuint32`, was previously `Uint64`.
  - `ActionTrace.ActionOrdinal` is now a `Varuint32`, was previously `uint32`.
  - `ActionTrace.CreatorActionOrdinal` is now a `Varuint32`, was previously `uint32`.
  - `ActionTrace.ClosestUnnotifiedAncestorActionOrdinal` is now a `Varuint32`, was previously `uint32`.
  - `Except.Code` is now a `Int64`, was previously `int`.
  - `ExceptLogContext.Level` is now a `ExceptLogLevel`, was previously `string`.
  - `ExceptLogContext.Line` is now a `uint64`, was previously `int`.

    **Note** While those are flagged as breaking change to augment the visibility, they are really bug fixes to fit with the behavior of `nodeos` directly.

- **BREAKING**: The decoding for ABI `variant` was not returning the correct `json` representation. Now ABI `variant` is returned as a two elements array, the first element being the `variant` type name as a `string` and the second the actual value as JSON. For example, assuming a `variant` type defined as `game_type: ["string", "uint32"]`, and a `field` of type `game_type`, before, the JSON serialization would have looked like `{"field":"some_string"}` or `{"field":100}` while after the change, it will get serialized to the correct form `{"field":["string", "some_string"]}` or `{"field":["uint32", 100]}`.

  **Note** While this is flagged as breaking change to augment the visibility, this is really a bug fix to fit with the behavior of `nodeos` directly.

- **BREAKING**: The serialization for `ExtendedAsset` was aligned with the `eos` codebase.  Beforehand, it would serialize the field name `"Contract"` with a capital `C`, and the `Asset` field as `"asset"` where it should have been `"quantity"`.

- **BREAKING**: Renamed `ConsoleLog` to `SafeString` for better re-usability in the codebase.

#### Added

- Proper handling for float precision in binary encoding/decoding.
- Added SHiP binary structures
- Added capabilities to read EOSIO Snapshot format (early implementation)
- Added architecture to support binary decoding/encoding Variant objects, see []
- Greatly improved performance of `NameToString` (`~230%`) method.
- `TimePoint` will decode with `0` nanoseconds, when the `fitNodeos` flag is set on the ABI.
- Ability to decode a `int128` and `uint128` in decimal format when `fitNodeos` flag is set on the ABI
- Ability to decode nested `arrays` in ABI decoder.
- Added `BlockState.Header` field of type `SignedBlockHeader` that was previously missing from the struct definition.
- Added `BlockState.AdditionalSignatures` field of type `[]ecc.Signature` that was previously missing from the struct definition.
- Added `ActionTrace.ContextFree` field of type `bool` that was previously missing from the struct definition.
- Normalized all logging to use `streamingfast/logging` and its trace enabled support.

#### Changed

- All errors are wrapped using `fmt.Errorf("...: %w", ..., err)` which is standard now in Go.

#### Fixed
- Optional encoding of primitive types.

  A struct with a non-pointer type tagged with `eos:"optional"` is now properly encoded at the binary level. **Important** that means that for non-pointer type, when the value of the type is the "emtpy" value according to Golang rules, it will be written as not-present at the binary level. If it's something that you do want want, use a pointer to a primitive type. It's actually a good habit to use a pointer type for "optional" element anyway, to increase awarness.

- Fix json tags for delegatebw action data.
- Unpacking binary `Except` now correctly works.
- Unpacking binary `Action` and `ActionTrace` now correctly works.
- Unpacking binary `TransactionTrace` now correctly works.
- Unpacking binary `TransactionReceipt` type will now correctly set the inner `TransactionWithID.ID` field correctly.
- Unpacking binary `BlockState` now correctly works but is restricted to EOSIO 2.0.x version.

#### Deprecated
- Renamed `AccountRAMDelta` to `AccountDelta` which is the correct name in EOSIO.
- Renamed `JSONFloat64` to `Float64`, to follow the same convention that was changed years ago with `Uint64`, etc. Type alias left for backwards compatibility.
