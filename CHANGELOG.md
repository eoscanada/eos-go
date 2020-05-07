# Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

* BREAKING: The serialization for `ExtendedAsset` was aligned with the `eos` codebase.  Beforehand, it would serialize the field name `"Contract"` with a capital `C`, and the `Asset` field as `"asset"` where it should have been `"quantity"`.

### Fixed

* Fix json tags for delegatebw action data.
