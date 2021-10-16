## Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

### Unreleased

#### Breaking Changes

#### Added

#### Changed

#### Fixed

#### Deprecated

### [0.10] (October 16th, 2021)

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
