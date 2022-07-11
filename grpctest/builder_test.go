package grpctest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCollections(t *testing.T) {
	items, err := NewCollections("./testdata/grpc_collection.json", nil)
	require.Nil(t, err)
	fmt.Println(items)
}
