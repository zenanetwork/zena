package integration

import (
	"testing"

	"github.com/zenanetwork/zena/tests/integration/indexer"
)

func TestKVIndexer(t *testing.T) {
	indexer.TestKVIndexer(t, CreateEvmd)
}
