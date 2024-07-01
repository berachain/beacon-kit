package serializer

import (
	"encoding/binary"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
)

func SerializeEnumerable(value types.SSZEnumerable[types.BaseSSZType]) ([]byte, error) {
	fixedParts := make([][]byte, 0)
	variableParts := make([][]byte, 0)

	// fixed_parts = [serialize(element) if not is_variable_size(element) else None for element in value]
	for _, e := range value.Elements() {
		if !e.IsFixed() {
			continue
		}

		bz, err := e.MarshalSSZ()
		if err != nil {
			return nil, err
		}

		// Append the serialized element to the fixed parts.
		fixedParts = append(fixedParts, bz)
	}

	// variable_parts = [serialize(element) if is_variable_size(element) else b"" for element in value]
	for _, e := range value.Elements() {
		if e.IsFixed() {
			continue
		}

		// Serialize the variable size element.
		bz, err := e.MarshalSSZ()
		if err != nil {
			return nil, err
		}

		// Append the serialized element to the variable parts.
		variableParts = append(variableParts, bz)
	}

	fixedSum := 0
	for _, part := range fixedParts {
		fixedSum += len(part)
	}

	// variable_lengths = [len(part) for part in variable_parts]
	variableLengths := make([]int, 0)
	variableSum := 0
	for _, part := range variableParts {
		variableLengths = append(variableLengths, len(part))
		variableSum += len(part)
	}

	if fixedSum+variableSum > 1<<(constants.BytesPerLengthOffset*constants.BitsPerByte) {
		return nil, errors.New("total length exceeds 2^64")
	}

	// Interleave offsets of variable-size parts with fixed-size parts
	variableOffsets := make([][]byte, 0)
	offset := uint32(fixedSum)
	for i := range value.Elements() {
		if !value.Elements()[i].IsFixed() {
			offsetBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(offsetBytes, offset)
			variableOffsets = append(variableOffsets, offsetBytes)
			offset += uint32(variableLengths[len(variableOffsets)-1])
		} else {
			variableOffsets = append(variableOffsets, nil)
		}
	}

	// Combine fixed parts with variable offsets
	combinedParts := make([][]byte, len(fixedParts))
	for i := range fixedParts {
		if variableOffsets[i] != nil {
			combinedParts[i] = variableOffsets[i]
		} else {
			combinedParts[i] = fixedParts[i]
		}
	}

	// Concatenate all parts
	result := make([]byte, 0, fixedSum+variableSum)
	for _, part := range combinedParts {
		result = append(result, part...)
	}
	for _, part := range variableParts {
		result = append(result, part...)
	}

	return result, nil

}
