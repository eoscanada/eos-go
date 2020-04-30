# Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Ability to decode nested `arrays`.

### Changed
- **BREAKING**: The decoding for ABI `variant` was not returning the correct `json` representation. Now ABI `variant` is returned as a two elements array, the first element being the `variant` type name as a `string` and the second the actual value as JSON. For example, assuming a `variant` type defined as `game_type: ["string", "uint32"]`, and a `field` of type `game_type`, before, the JSON serialization would have looked like `{"field":"some_string"}` or `{"field":100}` while after the change, it will get serialized to the correct form `{"field":["string", "some_string"]}` or `{"field":["uint32", 100]}`.

  **Note** While this is flagged a breaking change to augment the visibility, this is really a bug fix to fit with the behavior of `nodeos` directly.

- **BREAKING**: The serialization for `ExtendedAsset` was aligned with the `eos` codebase.  Beforehand, it would serialize the field name `"Contract"` with a capital `C`, and the `Asset` field as `"asset"` where it should have been `"quantity"`.

- **BREAKING**: We started adding an initial `context.Context` to all interruptible functions. All method performing an HTTP call have the new parameter as well as a bunch of other method. We cannot list all of them. If the caller already have a `context.Context` value, pass it to the function that now require one. Otherwise, simply pass `context.Background()`.
