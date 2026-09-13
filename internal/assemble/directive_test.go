package assemble

import (
	"maps"
	"testing"
)

func TestParseDirective(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		ok      bool
		want    directive
		wantErr bool
	}{
		{name: "bare", line: "::BASE", ok: true, want: directive{name: "BASE"}},
		{name: "surrounding spaces", line: "  ::BASE  ", ok: true, want: directive{name: "BASE"}},
		{
			name: "attrs",
			line: `::install{pkg="core" title="Two Words"}`,
			ok:   true,
			want: directive{name: "install", attrs: map[string]string{"pkg": "core", "title": "Two Words"}},
		},
		{
			name: "label and attrs",
			line: `::badge[Build Status]{color="green"}`,
			ok:   true,
			want: directive{name: "badge", label: "Build Status", hasLabel: true, attrs: map[string]string{"color": "green"}},
		},
		{
			name: "braces inside quoted value",
			line: `::install{pkg="{{ pkg }}"}`,
			ok:   true,
			want: directive{name: "install", attrs: map[string]string{"pkg": "{{ pkg }}"}},
		},
		{name: "heading", line: "# Title"},
		{name: "mid sentence", line: "see ::BASE here"},
		{name: "colons only", line: "::"},
		{name: "empty", line: ""},
		{name: "unquoted value", line: "::BASE{pkg=core}", ok: true, wantErr: true},
		{name: "attrs not separated", line: `::BASE{a="1"b="2"}`, ok: true, wantErr: true},
		{name: "duplicate attr", line: `::BASE{a="1" a="2"}`, ok: true, wantErr: true},
		{name: "trailing text", line: "::BASE and more", ok: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := parseDirective(tt.line)

			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.name != tt.want.name || got.label != tt.want.label || got.hasLabel != tt.want.hasLabel {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
			if len(got.attrs) != len(tt.want.attrs) || (len(tt.want.attrs) > 0 && !maps.Equal(got.attrs, tt.want.attrs)) {
				t.Errorf("attrs = %v, want %v", got.attrs, tt.want.attrs)
			}
		})
	}
}
