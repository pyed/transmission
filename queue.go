package transmission

import "context"

type queueArgs struct {
	IDs []int `json:"ids"`
}

// QueueMoveTop moves the specified torrents to the top of the queue.
func (c *Client) QueueMoveTop(ctx context.Context, ids ...int) error {
	return c.Execute(ctx, "queue-move-top", queueArgs{IDs: ids}, nil)
}

// QueueMoveUp moves the specified torrents up by one position in the queue.
func (c *Client) QueueMoveUp(ctx context.Context, ids ...int) error {
	return c.Execute(ctx, "queue-move-up", queueArgs{IDs: ids}, nil)
}

// QueueMoveDown moves the specified torrents down by one position in the queue.
func (c *Client) QueueMoveDown(ctx context.Context, ids ...int) error {
	return c.Execute(ctx, "queue-move-down", queueArgs{IDs: ids}, nil)
}

// QueueMoveBottom moves the specified torrents to the bottom of the queue.
func (c *Client) QueueMoveBottom(ctx context.Context, ids ...int) error {
	return c.Execute(ctx, "queue-move-bottom", queueArgs{IDs: ids}, nil)
}
