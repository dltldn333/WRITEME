package doc

import "testing"

func TestSplit(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    Document
		wantErr bool
	}{
		{
			name:   "no frontmatter",
			source: "# Hello\n",
			want:   Document{Frontmatter: "", Body: "# Hello\n"},
		},
		{
			name:   "frontmatter and body",
			source: "---\nprops:\n  - title\n---\n\n# Hello\n",
			want:   Document{Frontmatter: "props:\n  - title\n", Body: "\n# Hello\n"},
		},
		{
			name:   "empty frontmatter",
			source: "---\n---\nbody\n",
			want:   Document{Frontmatter: "", Body: "body\n"},
		},
		{
			name:   "horizontal rule is not frontmatter",
			source: "# Title\n\n---\n\nmore\n",
			want:   Document{Frontmatter: "", Body: "# Title\n\n---\n\nmore\n"},
		},
		{
			name:   "empty input",
			source: "",
			want:   Document{Frontmatter: "", Body: ""},
		},
		{
			name:    "unterminated frontmatter",
			source:  "---\nprops:\n  - title\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Split(tt.source)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Frontmatter != tt.want.Frontmatter {
				t.Errorf("Frontmatter\n got: %q\nwant: %q", got.Frontmatter, tt.want.Frontmatter)
			}
			if got.Body != tt.want.Body {
				t.Errorf("Body\n got: %q\nwant: %q", got.Body, tt.want.Body)
			}
		})
	}
}
