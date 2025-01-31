package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func MakeTempHomeDir(t *testing.T) string {
	t.Helper()
	// create random suffix to avoid conflicts
	const rndSuffixLen = 5
	bytes := make([]byte, rndSuffixLen)
	_, err := rand.Read(bytes)
	require.NoError(t, err)

	rndSuffix := hex.EncodeToString(bytes)

	tmpFolder := filepath.Join(os.TempDir(), "/injected-consensus", rndSuffix)
	require.NoError(t, os.MkdirAll(tmpFolder, os.ModePerm))
	return tmpFolder
}

func DeleteTempHomeDir(t *testing.T, homedir string) {
	t.Helper()
	err := os.RemoveAll(homedir)
	require.NoError(t, err)
}
