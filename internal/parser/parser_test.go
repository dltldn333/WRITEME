package parser

import "testing"

func TestParseDirective(t *testing.T) {
	tests := []struct {
		name string
		line string
		ok   bool
		want Directive
	}{
		{
			name: "bare name",
			line: "::BASE",
			ok:   true,
			want: Directive{Name: "BASE", Raw: "::BASE"},
		},
		{
			name: "attrs only",
			line: `::install{pkgName="core"}`,
			ok:   true,
			want: Directive{
				Name:  "install",
				Attrs: map[string]string{"pkgName": "core"},
				Raw:   `::install{pkgName="core"}`,
			},
		},
		{
			name: "label and attrs",
			line: `::badge[Build Status]{color="green"}`,
			ok:   true,
			want: Directive{
				Name:  "badge",
				Label: "Build Status",
				Attrs: map[string]string{"color": "green"},
				Raw:   `::badge[Build Status]{color="green"}`,
			},
		},
		{
			name: "multiple attrs",
			line: `::include{src="./install.md" title="Core"}`,
			ok:   true,
			want: Directive{
				Name:  "include",
				Attrs: map[string]string{"src": "./install.md", "title": "Core"},
				Raw:   `::include{src="./install.md" title="Core"}`,
			},
		},
		{
			name: "label only",
			line: "::note[heads up]",
			ok:   true,
			want: Directive{Name: "note", Label: "heads up", Raw: "::note[heads up]"},
		},
		{name: "plain heading", line: "# Title"},
		{name: "not at line start", line: "see ::badge for details"},
		{name: "single colon", line: ":badge"},
		{name: "empty line", line: ""},
		{name: "colons only", line: "::"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseDirective(tt.line)

			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if !tt.ok {
				return
			}
			if got.Name != tt.want.Name {
				t.Errorf("Name = %q, want %q", got.Name, tt.want.Name)
			}
			if got.Label != tt.want.Label {
				t.Errorf("Label = %q, want %q", got.Label, tt.want.Label)
			}
			if got.Raw != tt.want.Raw {
				t.Errorf("Raw = %q, want %q", got.Raw, tt.want.Raw)
			}
			if len(got.Attrs) != len(tt.want.Attrs) {
				t.Fatalf("Attrs = %v, want %v", got.Attrs, tt.want.Attrs)
			}
			for k, v := range tt.want.Attrs {
				if got.Attrs[k] != v {
					t.Errorf("Attrs[%q] = %q, want %q", k, got.Attrs[k], v)
				}
			}
		})
	}
}
