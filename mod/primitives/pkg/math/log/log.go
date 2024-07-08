package log

import "math/bits"

// ILog2Ceil returns the ceiling of the base 2 logarithm of the input.
func ILog2Ceil[U64T ~uint64](u U64T) uint8 {
	// Log2(0) is undefined, should we panic?
	if u == 0 {
		return 0
	}
	//#nosec:G701 // we handle the case of u == 0 above, so this is safe.
	return uint8(bits.Len64(uint64(u - 1)))
}

// ILog2Floor returns the floor of the base 2 logarithm of the input.
func ILog2Floor[U64T ~uint64](u U64T) uint8 {
	// Log2(0) is undefined, should we panic?
	if u == 0 {
		return 0
	}
	//#nosec:G701 // we handle the case of u == 0 above, so this is safe.
	return uint8(bits.Len64(uint64(u))) - 1
}
