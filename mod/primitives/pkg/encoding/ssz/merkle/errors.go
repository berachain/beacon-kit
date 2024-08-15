package merkle

import "errors"

// ErrUnexpectedProofLength is returned when the proof length is unexpected.
var ErrUnexpectedProofLength = errors.New("unexpected proof length")

// ErrMistmatchLeavesIndicesLength is returned when the leaves and indices length mismatch.
var ErrMistmatchLeavesIndicesLength = errors.New("mismatched leaves and indices length")
