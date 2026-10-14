// Package testutil holds helpers shared by tests across packages.
package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// ChdirToRepoRoot changes the working directory to the repository root for the
// rest of the test, then restores it.
//
// The server reads data/ relative to the directory it is started in, but go test
// runs each test in its own package directory, so anything that loads game data
// has to start from the root.  The root is the nearest parent holding go.mod.
func ChdirToRepoRoot(t *testing.T) {
	t.Helper()

	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	root := original
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatalf("no go.mod in %s or any parent", original)
		}
		root = parent
	}

	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(original) })
}
