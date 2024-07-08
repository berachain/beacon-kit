package pow

// NextPowerOfTwo returns the next power of 2 for the given input.
//
//nolint:mnd // todo fix.
func PrevPowerOfTwo[U64T ~uint64](u U64T) U64T {
	if u == 0 {
		return 1
	}
	u |= u >> 1
	u |= u >> 2
	u |= u >> 4
	u |= u >> 8
	u |= u >> 16
	u |= u >> 32
	return u - (u >> 1)
}

// NextPowerOfTwo returns the next power of 2 for the given input.
//
//nolint:mnd // todo fix.
func NextPowerOfTwo[U64T ~uint64](u U64T) U64T {
	if u == 0 {
		return 1
	}
	if u > 1<<63 {
		panic("Next power of 2 is 1 << 64.")
	}
	u--
	u |= u >> 1
	u |= u >> 2
	u |= u >> 4
	u |= u >> 8
	u |= u >> 16
	u |= u >> 32
	u++
	return u
}
