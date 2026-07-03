package util_test

import (
	"testing"

	"github.com/marcuwynu23/git-share/internal/util"
)

func TestFatalExists(t *testing.T) {
	// Fatal calls os.Exit(1) so we can't directly test it.
	// Just verify the function value is not nil.
	t.Log("Fatal function exists (cannot directly test os.Exit)")
	_ = util.Fatal
}
