package collections_test

import (
	"testing"

	"github.com/itsdevbear/bolaris/lib/store/collections"

	sdk "cosmossdk.io/collections"
	sdkcollections "cosmossdk.io/collections"
	"github.com/stretchr/testify/require"
)

func TestCircularQueuePushPop(t *testing.T) {
	sk, ctx := deps()
	schema := sdk.NewSchemaBuilder(sk)
	queue := collections.NewCircularQueue[uint64](schema, "testQueue", sdkcollections.Uint64Value, 5)

	// Push elements into the queue
	for i := uint64(0); i < 5; i++ {
		_, err := queue.Push(ctx, i)
		require.NoError(t, err)
	}

	// Push another element, which should cause the first element to be evicted (circular behavior)
	evicted, err := queue.Push(ctx, 5)
	require.NoError(t, err)
	require.Equal(t, uint64(0), evicted)

	// Pop elements and verify order
	for i := uint64(1); i <= 5; i++ {
		item, err := queue.Peek(ctx)
		require.NoError(t, err)
		require.Equal(t, i, item)
		_, err = queue.Push(ctx, item+5) // Push next element to maintain circularity
		require.NoError(t, err)
	}
}

func TestCircularQueueOverflow(t *testing.T) {
	sk, ctx := deps()
	schema := sdk.NewSchemaBuilder(sk)
	queue := collections.NewCircularQueue[uint64](schema, "testQueueOverflow", sdkcollections.Uint64Value, 3)

	// Push elements more than the queue size to test overflow
	for i := uint64(0); i < 10; i++ {
		_, err := queue.Push(ctx, i)
		require.NoError(t, err)
	}

	// Verify that the queue only contains the last 3 elements
	expected := []uint64{7, 8, 9}
	for _, e := range expected {
		item, err := queue.Peek(ctx)
		require.NoError(t, err)
		require.Equal(t, e, item)
		_, err = queue.Push(ctx, item) // Push the peeked item back to simulate continuous operation
		require.NoError(t, err)
	}
}
