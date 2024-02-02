package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dememorized/forsikt/auditfile"
	"github.com/dememorized/forsikt/graph"
	"golang.org/x/mod/semver"
)

func main() {
	mods, err := graph.GraphModule(context.Background())
	if err != nil {
		panic(err)
	}

	filename := "go.audit"
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	af, err := auditfile.Parse(filename, bytes)
	if err != nil {
		panic(err)
	}

	trusts := map[string]map[string]string{}
	for _, t := range af.Trust {
		if t == nil {
			continue
		}
		if _, exists := trusts[t.Mod.Path]; !exists {
			trusts[t.Mod.Path] = make(map[string]string)
		}

		trusts[t.Mod.Path][t.Mod.High] = t.Mod.Low
	}

	violations := map[string]map[string]string{}
	for _, t := range af.Violation {
		if t == nil {
			continue
		}
		if _, exists := violations[t.Mod.Path]; !exists {
			violations[t.Mod.Path] = make(map[string]string)
		}

		violations[t.Mod.Path][t.Mod.High] = t.Mod.Low
	}

	resNotFounds := []string{}
	resVersionMismatch := []string{}
	resViolations := []string{}

	for _, mod := range mods.Versions {
		p, v := mod.Path, mod.Version
		if _, exists := violations[p]; exists {
			if versionInRanges(v, violations[p]) {
				resViolations = append(resViolations, mod.String())
				continue
			}
		}

		if _, exists := trusts[p]; exists {
			if versionInRanges(v, trusts[p]) {
				continue
			}
			resVersionMismatch = append(resVersionMismatch, mod.String())
		} else {
			resNotFounds = append(resNotFounds, mod.String())
		}
	}

	if len(resNotFounds) == 0 && len(resVersionMismatch) == 0 && len(resViolations) == 0 {
		// all good!
		return
	}

	fmt.Printf("modules missing from '%s':\n\t%s\n\n", filename, strings.Join(resNotFounds, "\n\t"))
	fmt.Printf("modules with wrong versions from '%s':\n\t%s\n\n", filename, strings.Join(resVersionMismatch, "\n\t"))
	fmt.Printf("policy violations from '%s':\n\t%s\n\n", filename, strings.Join(resViolations, "\n\t"))

	fmt.Printf("### audit failed ###\n%d missing modules, %d version mismatches, %d policy violations\n",
		len(resNotFounds), len(resVersionMismatch), len(resViolations))

	os.Exit(1)
}

func versionInRanges(v string, ranges map[string]string) bool {
	for upper, lower := range ranges {
		// catch-all '*'
		if upper == "" {
			return true
		}

		if semver.Compare(v, upper) <= 0 && semver.Compare(v, lower) >= 0 {
			return true
		}
	}
	return false
}
