// Copyright 2022 Emil Tullstedt.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package auditfile

import (
	"fmt"
	"strings"
)

type LexError struct {
	Filename string
	Pos      Position
	Err      error
}

func (e LexError) Error() string {
	return fmt.Sprintf("%s:%s: %s", e.Filename, e.Pos, e.Err.Error())
}

type ParserError struct {
	Filename string
	Pos      Position
	Err      error
	Verb     string
	ModPath  string
}

func (e ParserError) Error() string {
	path := e.Filename + ":"
	if e.Pos.Line != 0 {
		path = fmt.Sprintf("%s:%d:", e.Filename, e.Pos.Line)
	}

	msg := []string{path}

	if e.ModPath != "" {
		msg = append(msg, e.ModPath)
	}
	if e.Verb != "" {
		msg = append(msg, e.Verb)
	}
	if e.Err != nil {
		msg = append(msg, e.Err.Error())
	}
	return strings.Join(msg, " ")
}

type ErrorList []error

func (e ErrorList) Error() string {
	errStrs := make([]string, len(e))
	for i, err := range e {
		errStrs[i] = err.Error()
	}
	return strings.Join(errStrs, "\n")
}
