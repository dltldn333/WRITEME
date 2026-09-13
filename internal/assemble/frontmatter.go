package assemble

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type frontmatter struct {
	Props []string `yaml:"props"`
}

// splitFrontmatter separates a leading --- block from the body. offset is how
// many lines the block used, so body line numbers can point into the original file.
func splitFrontmatter(src string) (fm frontmatter, body string, offset int, err error) {
	if !strings.HasPrefix(src, "---\n") {
		return frontmatter{}, src, 0, nil
	}

	rest := src[len("---\n"):]
	var block string

	switch {
	case strings.HasPrefix(rest, "---\n"):
		body = rest[len("---\n"):]
	case rest == "---":
		body = ""
	default:
		if i := strings.Index(rest, "\n---\n"); i >= 0 {
			block, body = rest[:i+1], rest[i+len("\n---\n"):]
		} else if strings.HasSuffix(rest, "\n---") {
			block, body = rest[:len(rest)-len("---")], ""
		} else {
			return frontmatter{}, "", 0, fmt.Errorf("frontmatter is not closed with ---")
		}
	}

	if err := yaml.Unmarshal([]byte(block), &fm); err != nil {
		return frontmatter{}, "", 0, fmt.Errorf("frontmatter: %w", err)
	}

	offset = strings.Count(src[:len(src)-len(body)], "\n")
	return fm, body, offset, nil
}
