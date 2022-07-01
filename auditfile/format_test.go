package auditfile

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	files, err := os.ReadDir("testdata")
	require.NoError(t, err)

	for _, entry := range files {
		if !strings.HasSuffix(entry.Name(), ".raw") {
			continue
		}

		bytes, err := os.ReadFile("testdata/" + entry.Name())
		require.NoError(t, err)

		f, err := Parse("go.audit", bytes)
		require.NoError(t, err)

		expected, err := os.ReadFile("testdata/" + strings.TrimSuffix(entry.Name(), ".raw") + ".gold")
		require.NoError(t, err)

		assert.Equal(t, string(expected), f.String())
	}
}
