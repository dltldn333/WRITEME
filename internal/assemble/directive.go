package assemble

import (
	"fmt"
	"regexp"
	"strings"
)

type directive struct {
	name     string
	label    string
	hasLabel bool
	attrs    map[string]string
}

var (
	directiveStart = regexp.MustCompile(`^::[A-Za-z]`)
	directiveRe    = regexp.MustCompile(`^::([A-Za-z][A-Za-z0-9_-]*)(\[[^\]]*\])?(\{(?:[^}"]|"[^"]*")*\})?$`)
	attrRe         = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_-]*)="([^"]*)"`)
)

// parseDirective reports ok=false for ordinary lines. A line that starts like a
// component but does not parse is an error, so typos never leak into README.md.
func parseDirective(line string) (d directive, ok bool, err error) {
	line = strings.TrimSpace(line)
	if !directiveStart.MatchString(line) {
		return directive{}, false, nil
	}

	m := directiveRe.FindStringSubmatch(line)
	if m == nil {
		return directive{}, true, fmt.Errorf("malformed component %q", line)
	}

	d.name = m[1]
	if m[2] != "" {
		d.hasLabel = true
		d.label = m[2][1 : len(m[2])-1]
	}
	if m[3] != "" {
		d.attrs, err = parseAttrs(m[3][1 : len(m[3])-1])
		if err != nil {
			return directive{}, true, fmt.Errorf("component %q: %w", line, err)
		}
	}
	return d, true, nil
}

func parseAttrs(s string) (map[string]string, error) {
	attrs := make(map[string]string)

	for s = strings.TrimSpace(s); s != ""; s = strings.TrimSpace(s) {
		m := attrRe.FindStringSubmatch(s)
		if m == nil {
			return nil, fmt.Errorf("bad attribute near %q (use key=\"value\")", s)
		}
		if _, dup := attrs[m[1]]; dup {
			return nil, fmt.Errorf("attribute %q is given twice", m[1])
		}
		attrs[m[1]] = m[2]

		s = s[len(m[0]):]
		if s != "" && s[0] != ' ' && s[0] != '\t' {
			return nil, fmt.Errorf("attributes must be separated by spaces near %q", s)
		}
	}
	return attrs, nil
}
