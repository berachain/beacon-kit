package cache

import "errors"

var (
	errInvalidIterator              = errors.New("invalid iterator")
	errInvalidIteratorRangeCacheKey = errors.New("invalid iterator range cache key")
)
