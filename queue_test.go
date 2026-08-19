package transmission

import (
	"context"
	"testing"
)

func TestQueueMovement(t *testing.T) {
	server, _ := newMockServer(t)
	defer server.Close()

	client, err := New(server.URL, "admin", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	if err := client.QueueMoveTop(ctx, 1, 2); err != nil {
		t.Errorf("QueueMoveTop failed: %v", err)
	}
	if err := client.QueueMoveUp(ctx, 1); err != nil {
		t.Errorf("QueueMoveUp failed: %v", err)
	}
	if err := client.QueueMoveDown(ctx, 1); err != nil {
		t.Errorf("QueueMoveDown failed: %v", err)
	}
	if err := client.QueueMoveBottom(ctx, 1, 2); err != nil {
		t.Errorf("QueueMoveBottom failed: %v", err)
	}
}
