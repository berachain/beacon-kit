package uint256_test

// func TestLittleEndian_UInt256(t *testing.T) {
// 	le := uint256.LittleEndian([]byte{1, 2, 3, 4, 5})
// 	expected := new(holimanuint256.Int).SetBytes([]byte{1, 2, 3, 4, 5})
// 	assert.Equal(t, expected, le.UInt256())
// }

// func TestLittleEndian_Big(t *testing.T) {
// 	le := uint256.LittleEndian([]byte{1, 2, 3, 4, 5})
// 	expected := new(holimanuint256.Int).SetBytes([]byte{1, 2, 3, 4, 5})
// 	assert.Equal(t, expected, le.Big())
// }

// func TestLittleEndian_MarshalJSON(t *testing.T) {
// 	le := uint256.LittleEndian([]byte{1, 2, 3, 4, 5})
// 	expected := []byte("\"0x0504030201\"")
// 	result, err := le.MarshalJSON()
// 	assert.NoError(t, err)
// 	assert.Equal(t, expected, result)
// }

// func TestLittleEndian_UnmarshalJSON(t *testing.T) {
// 	le := new(uint256.LittleEndian)
// 	err := le.UnmarshalJSON([]byte("0x0504030201"))
// 	assert.NoError(t, err)
// 	assert.Equal(t, uint256.LittleEndian([]byte{1, 2, 3, 4, 5}), *le)
// }
