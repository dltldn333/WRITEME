package assemble

import (
	"slices"
	"testing"
)

func TestSplitFrontmatter(t *testing.T) {
	tests := []struct {
		name       string
		src        string
		wantProps  []string
		wantBody   string
		wantOffset int
		wantErr    bool
	}{
		{name: "no frontmatter", src: "# Hello\n", wantBody: "# Hello\n"},
		{name: "empty input", src: ""},
		{
			name:       "props and body",
			src:        "---\nprops:\n  - title\n  - version\n---\n\n# Hello\n",
			wantProps:  []string{"title", "version"},
			wantBody:   "\n# Hello\n",
			wantOffset: 5,
		},
		{name: "empty block", src: "---\n---\nbody\n", wantBody: "body\n", wantOffset: 2},
		{name: "closed at end of file", src: "---\nprops: [a]\n---", wantProps: []string{"a"}, wantOffset: 2},
		{name: "horizontal rule is not frontmatter", src: "# T\n\n---\n\nmore\n", wantBody: "# T\n\n---\n\nmore\n"},
		{name: "not closed", src: "---\nprops: [a]\n", wantErr: true},
		{name: "invalid yaml", src: "---\nprops: [a\n---\n", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body, offset, err := splitFrontmatter(tt.src)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !slices.Equal(fm.Props, tt.wantProps) {
				t.Errorf("props = %v, want %v", fm.Props, tt.wantProps)
			}
			if body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
			if offset != tt.wantOffset {
				t.Errorf("offset = %d, want %d", offset, tt.wantOffset)
			}
		})
	}
}
