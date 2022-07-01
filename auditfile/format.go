package auditfile

import (
	"fmt"
	"golang.org/x/mod/semver"
	"sort"
	"strings"
)

func (f File) String() string {
	var lines []string
	if f.Audit != nil {
		lines = append(lines, fmt.Sprintf("audit %s\n", f.Audit.Version))
	}

	SortTrust(f.Trust)
	switch len(f.Trust) {
	case 0:
	// no-op
	case 1:
		lines = append(lines, f.Trust[0].String())
	default:
		lines = append(lines, "trust (")
		for _, trust := range f.Trust {
			lines = append(lines, trust.stringWithVerb("", true))
		}
		lines = append(lines, ")\n")
	}

	SortViolation(f.Violation)
	switch len(f.Violation) {
	case 0:
	// no-op
	case 1:
		lines = append(lines, f.Violation[0].String())
	default:
		lines = append(lines, "violation (")
		for _, violation := range f.Violation {
			lines = append(lines, Trust(*violation).stringWithVerb("", true))
		}
		lines = append(lines, ")\n")
	}

	return strings.Join(lines, "\n") + "\n"
}

func SortTrust(list []*Trust) {
	sort.Slice(list, func(i, j int) bool {
		a, b := list[i], list[j]
		if a.Mod.Path == b.Mod.Path {
			comp := semver.Compare(a.Mod.High, b.Mod.High)
			return comp == -1
		}
		return a.Mod.Path < b.Mod.Path
	})
}

func SortViolation(list []*Violation) {
	sort.Slice(list, func(i, j int) bool {
		a, b := list[i], list[j]
		if a.Mod.Path == b.Mod.Path {
			comp := semver.Compare(a.Mod.High, b.Mod.High)
			return comp == -1
		}
		return a.Mod.Path < b.Mod.Path
	})
}
