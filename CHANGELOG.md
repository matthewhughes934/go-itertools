# Changelog

## 0.5.0 - 2025-06-20

### Added

  - Add `Max`, `MaxFunc`, `Min`, and `MinFunc` to find the min or max values
  - Add `Pairwise` to iterate pairwise over a sequence
  - Add `SliceFrom` to slice from a position within the sequence
  - Add `Flatten` and `Flatten` to iterate across keys then values of a
    `iter.Seq2`

### Changed

  - `Slice` now accepts a negative number for 'end', meaning it will continue
    until it exhausts the sequence

## 0.4.0 - 2024-10-28

### Added

  - Add `ChainSlices` as a convenience for `Chain`ing over slices
  - Add `ChainMaps` as a convenience for `Chain`ing over maps

## 0.3.0 - 2024-10-05

### Added

  - Add functions to return first matching elements

## 0.2.0 - 2024-09-05

### Fixed

  - Fix race condition in `IterCtx*` functions

## 0.1.0 - 2024-08-21

*Initial release*
