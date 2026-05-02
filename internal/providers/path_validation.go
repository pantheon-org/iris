package providers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pantheon-org/iris/internal/ierrors"
)

// ValidateProjectRoot checks that projectRoot does not contain path traversal
// sequences. An empty string is allowed — callers interpret it as "use global
// config path". Returns ErrPathTraversal if any ".." component is present.
func ValidateProjectRoot(projectRoot string) error {
	if projectRoot == "" {
		return nil
	}
	// Normalise separators then split on '/' to inspect every path component.
	for _, part := range strings.FieldsFunc(filepath.ToSlash(projectRoot), func(r rune) bool { return r == '/' }) {
		if part == ".." {
			return fmt.Errorf("projectRoot %q: %w", projectRoot, ierrors.ErrPathTraversal)
		}
	}
	return nil
}
