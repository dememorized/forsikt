package auditfile

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		file   *File
		errStr string
	}{
		{
			name:   "blank file",
			in:     "",
			errStr: "testfile.audit: expected 'audit' directive",
		},
		{
			name: "base file",
			in:   "audit 1",
			file: &File{
				Audit: &FileVersion{
					Version: "1",
				},
				Trust:     nil,
				Violation: nil,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			f, err := Parse("testfile.audit", []byte(tc.in))
			if err != nil {
				if strings.TrimSpace(err.Error()) != tc.errStr {
					t.Errorf("got unexpected error: %v", err)
				}
			}

			if tc.file == nil {
				if f != nil {
					t.Errorf("expected <nil> when parsing file, got %#v", f)
				}
				return
			}

			if tc.file.Audit != nil {
				if f.Audit == nil {
					t.Errorf("expected 'audit' directive, got <nil>")
				}

				if f.Audit.Version != tc.file.Audit.Version {
					t.Errorf("expected 'audit' version %s, got %s", tc.file.Audit.Version, f.Audit.Version)
				}
			}

			if f.Trust != nil || f.Violation != nil {
				t.Error("testing trust/violation is not yet implemented")
				return
			}
		})
	}
}
