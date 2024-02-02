// Copyright 2022 Emil Tullstedt.
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auditfile

import (
	"errors"
	"fmt"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	"strconv"
	"strings"
	"unicode"
)

type File struct {
	Audit     *FileVersion
	Trust     []*Trust
	Violation []*Violation

	Syntax *FileSyntax
}

type FileVersion struct {
	Version string
	Syntax  *Line
}

type VersionInterval struct {
	Path      string
	Low, High string
}

func (v VersionInterval) String() string {
	version := v.Low
	if semver.Compare(v.Low, v.High) != 0 {
		version = fmt.Sprintf("[%s, %s]", v.Low, v.High)
	}

	return v.Path + " " + version
}

type Trust struct {
	Mod       VersionInterval
	Notes     string
	Signature string
	Syntax    *Line
}

func (t Trust) String() string {
	return t.stringWithVerb("trust", false)
}

func (t Trust) stringWithVerb(verb string, inBlock bool) string {
	var lines []string
	var inlineComment string

	if t.Syntax != nil {
		for _, line := range t.Syntax.Comments.Before {
			lines = append(lines, line.Token)
		}

		inlineComments := []string{}
		for _, c := range t.Syntax.Comments.Suffix {
			inlineComments = append(inlineComments, c.Token)
		}
		inlineComment = strings.Join(inlineComments, " ")
		if inlineComment != "" {
			inlineComment = " " + inlineComment
		}
	}

	v := verb + " "
	if inBlock {
		v = ""
	}

	lines = append(lines, fmt.Sprintf("%s%s %s%s", v, t.Mod.String(), t.Signature, inlineComment))

	if inBlock {
		for i, line := range lines {
			lines[i] = "\t" + line
		}
	}

	return strings.Join(lines, "\n")
}

type Violation Trust

func (v Violation) String() string {
	return Trust(v).stringWithVerb("violation", false)
}

func Parse(file string, data []byte) (parsed *File, err error) {
	fs, err := parse(file, data)
	if err != nil {
		return nil, err
	}
	f := &File{
		Syntax: fs,
	}
	var errs ErrorList

	for _, x := range fs.Stmt {
		switch x := x.(type) {
		case *Line:
			f.add(&errs, x, x.Token[0], x.Token[1:])

		case *LineBlock:
			if len(x.Token) > 1 {
				errs = append(errs, ParserError{
					Filename: file,
					Pos:      x.Start,
					Err:      fmt.Errorf("unknown block type: %s", strings.Join(x.Token, " ")),
				})
				continue
			}
			switch x.Token[0] {
			default:
				errs = append(errs, ParserError{
					Filename: file,
					Pos:      x.Start,
					Err:      fmt.Errorf("unknown block type: %s", strings.Join(x.Token, " ")),
				})
				continue
			case "trust", "violation":
				for _, l := range x.Line {
					f.add(&errs, l, x.Token[0], l.Token)
				}
			}
		}
	}

	if f.Audit == nil {
		errs = append(errs, ParserError{
			Filename: file,
			Pos:      Position{},
			Err:      fmt.Errorf("expected 'audit' directive"),
		})
	} else if f.Audit.Version != "1" {
		errs = append(errs, ParserError{
			Filename: file,
			Pos:      f.Audit.Syntax.Start,
			Err:      fmt.Errorf("expected 'audit' directive version to be equal to '1'"),
		})
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return f, nil
}

func (f *File) add(errs *ErrorList, line *Line, verb string, args []string) {
	wrapError := func(err error) {
		*errs = append(*errs, ParserError{
			Filename: f.Syntax.Name,
			Pos:      line.Start,
			Err:      err,
		})
	}
	errorf := func(format string, args ...interface{}) {
		wrapError(fmt.Errorf(format, args...))
	}

	switch verb {
	default:
		errorf("unknown directive: %s", verb)

	case "audit":
		if f.Audit != nil {
			errorf("repeated audit statement")
			return
		}
		if len(args) != 1 {
			errorf("audit directive expects exactly one argument")
			return
		}

		f.Audit = &FileVersion{
			Version: args[0],
			Syntax:  line,
		}

	case "trust", "violation":
		if len(args) != 3 && len(args) != 7 {
			errorf("usage: %s module/path v1.2.3 \"Nils Holgersson <nils.holgersson@example.org>\"", verb)
			return
		}
		s, err := parseString(&args[0])
		if err != nil {
			errorf("invalid quoted string: %v", err)
			return
		}
		args = args[1:]
		v, err := parseVersionInterval(verb, s, &args)
		if err != nil {
			wrapError(err)
			return
		}

		if len(args) != 1 {
			errorf("missing signature")
			return
		}

		signature, err := parseString(&args[0])
		if err != nil {
			errorf("invalid quoted string: %v", err)
			return
		}

		notes := parseNotes(line)

		switch verb {
		case "trust":
			f.Trust = append(f.Trust, &Trust{
				Mod:       v,
				Notes:     notes,
				Signature: signature,
				Syntax:    line,
			})
		case "violation":
			f.Violation = append(f.Violation, &Violation{
				Mod:       v,
				Notes:     notes,
				Signature: signature,
				Syntax:    line,
			})
		}
	}
}

func parseNotes(line *Line) string {
	if line == nil {
		return ""
	}

    trim := func (s string) string {
        return strings.TrimSpace(strings.TrimPrefix("//", s))
    }

	var notes []string
	for _, note := range line.Comments.Before {
		notes = append(notes, trim(note.Token))
	}
    for _, note := range line.Comments.Suffix {
		notes = append(notes, trim(note.Token))
    }

	return strings.Join(notes, " ")
}

func parseVersion(verb string, path string, s *string) (string, error) {
	t, err := parseString(s)
	if err != nil {
		return "", ParserError{
			Verb:    verb,
			ModPath: path,
			Err: &module.InvalidVersionError{
				Version: *s,
				Err:     err,
			},
		}
	}

	cv := module.CanonicalVersion(t)
	if cv == "" && t != "*" {
		return "", ParserError{
			Verb:    verb,
			ModPath: path,
			Err: &module.InvalidVersionError{
				Version: t,
				Err:     errors.New("must be of the form v1.2.3"),
			},
		}
	}
	t = cv

	*s = t
	return *s, nil
}

func parseString(s *string) (string, error) {
	t := *s
	if strings.HasPrefix(t, `"`) {
		var err error
		if t, err = strconv.Unquote(t); err != nil {
			return "", err
		}
	} else if strings.ContainsAny(t, "\"'`") {
		// Other quotes are reserved both for possible future expansion
		// and to avoid confusion. For example if someone types 'x'
		// we want that to be a syntax error and not a literal x in literal quotation marks.
		return "", fmt.Errorf("unquoted string cannot contain quote")
	}
	return AutoQuote(t), nil
}

// AutoQuote returns s or, if quoting is required for s to appear in a go.mod,
// the quotation of s.
func AutoQuote(s string) string {
	if MustQuote(s) {
		return strconv.Quote(s)
	}
	return s
}

// MustQuote reports whether s must be quoted in order to appear as
// a single token in a go.mod line.
func MustQuote(s string) bool {
	for _, r := range s {
		switch r {
		case ' ', '"', '\'', '`':
			return true

		case '(', ')', '[', ']', '{', '}', ',':
			if len(s) > 1 {
				return true
			}

		default:
			if !unicode.IsPrint(r) {
				return true
			}
		}
	}
	return s == "" || strings.Contains(s, "//") || strings.Contains(s, "/*")
}

func parseVersionInterval(verb string, path string, args *[]string) (VersionInterval, error) {
	toks := *args
	if len(toks) == 0 || toks[0] == "(" {
		return VersionInterval{}, fmt.Errorf("expected '[' or version")
	}
	if toks[0] != "[" {
		v, err := parseVersion(verb, path, &toks[0])
		if err != nil {
			return VersionInterval{}, err
		}
		*args = toks[1:]
		return VersionInterval{Path: path, Low: v, High: v}, nil
	}
	toks = toks[1:]

	if len(toks) == 0 {
		return VersionInterval{}, fmt.Errorf("expected version after '['")
	}
	low, err := parseVersion(verb, path, &toks[0])
	if err != nil {
		return VersionInterval{}, err
	}
	toks = toks[1:]

	if len(toks) == 0 || toks[0] != "," {
		return VersionInterval{}, fmt.Errorf("expected ',' after version")
	}
	toks = toks[1:]

	if len(toks) == 0 {
		return VersionInterval{}, fmt.Errorf("expected version after ','")
	}
	high, err := parseVersion(verb, path, &toks[0])
	if err != nil {
		return VersionInterval{}, err
	}
	toks = toks[1:]

	if len(toks) == 0 || toks[0] != "]" {
		return VersionInterval{}, fmt.Errorf("expected ']' after version")
	}
	toks = toks[1:]

	*args = toks
	return VersionInterval{Path: path, Low: low, High: high}, nil
}
