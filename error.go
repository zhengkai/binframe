package binframe

import "errors"

// ErrSizeOverflow size overflow
var ErrSizeOverflow = errors.New(`size overflow (uint64)`)

// ErrOversized message size is bigger than SizeThreshold
var ErrOversized = errors.New(`oversized message`)
