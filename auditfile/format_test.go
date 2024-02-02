package auditfile

import (
	"os"
	"strings"
	"testing"

	"github.com/dememorized/forsikt/internal/hoppsan"
)

func TestFormat(t *testing.T) {
	files, err := os.ReadDir("testdata")
	hoppsan.NoError(t, err)

	for _, entry := range files {
		if !strings.HasSuffix(entry.Name(), ".raw") {
			continue
		}

		bytes, err := os.ReadFile("testdata/" + entry.Name())
		hoppsan.NoError(t, err)

		f, err := Parse("go.audit", bytes)
		hoppsan.NoError(t, err)

		expected, err := os.ReadFile("testdata/" + strings.TrimSuffix(entry.Name(), ".raw") + ".gold")
		hoppsan.NoError(t, err)

		hoppsan.Equal(t, string(expected), f.String())
	}
}
