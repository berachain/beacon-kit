package hex

// has0xPrefix returns true if s has a 0x prefix.
func has0xPrefix[T []byte | string](s T) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

// isQuotedString returns true if input has quotes.
func isQuotedString[T []byte | string](input T) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}
