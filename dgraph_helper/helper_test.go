package dgraph_helper

import (
	"os"
	"testing"
)

// *** Drop relative data in dgraph, or test results may be affected ***
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
