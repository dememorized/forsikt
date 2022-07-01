package graph

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	"os/exec"
	"sort"
	"strings"
)

type Modules struct {
	Versions []module.Version
}

func GraphModule(ctx context.Context) (*Modules, error) {
	versions := []module.Version{}

	cmd := exec.CommandContext(ctx, "go", "mod", "graph")
	buf := bytes.Buffer{}
	errMsg := bytes.Buffer{}
	cmd.Stdout = &buf
	cmd.Stderr = &errMsg

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error running 'go mod graph'\nerror: %w\nstderr: %s", err, errMsg.String())
	}

	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		modules := strings.Split(line, " ")
		if len(modules) != 2 {
			return nil, fmt.Errorf("expected line to have two modules, got: %s\n", line)
		}

		target := modules[1]
		parts := strings.Split(target, "@")
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected target to have version, got: %s\n", target)
		}

		versions = append(versions, module.Version{
			Path:    parts[0],
			Version: parts[1],
		})
	}

	sort.Slice(versions, func(i, j int) bool {
		a, b := versions[i], versions[j]
		if a.Path == b.Path {
			comp := semver.Compare(a.Version, b.Version)
			return comp == -1
		}
		return a.Path < b.Path
	})

	versionsMerged := []module.Version{}

	var prev module.Version
	for _, v := range versions {
		if v.Path == prev.Path && v.Version == prev.Version {
			continue
		}
		versionsMerged = append(versionsMerged, v)
		prev = v
	}

	return &Modules{Versions: versionsMerged}, nil
}
