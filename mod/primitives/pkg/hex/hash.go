package hex

import "encoding/hex"

// hex helpers for the Hash type.

// SetupHashFormat returns a hex string with 0x prefix.
func HexEncodeWithPrefix[B ~[]byte](h B) []byte {
	hexb := make([]byte, 2+len(h)*2)
	copy(hexb, "0x")
	hex.Encode(hexb[2:], h[:])
	return hexb
}
