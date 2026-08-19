package transmission

import (
	"context"
	"fmt"
)

type idsArgs struct {
	IDs []int `json:"ids,omitempty"`
}

func (c *Client) sendIDsCommand(ctx context.Context, method string, ids ...int) error {
	var args any
	if len(ids) > 0 {
		args = idsArgs{IDs: ids}
	}
	return c.Execute(ctx, method, args, nil)
}

// StartTorrents starts one or more torrents by ID.
func (c *Client) StartTorrents(ctx context.Context, ids ...int) error {
	return c.sendIDsCommand(ctx, "torrent-start", ids...)
}

// StartTorrent starts a single torrent by ID (backwards-compatible signature).
func (c *Client) StartTorrent(id int) (string, error) {
	err := c.StartTorrents(context.Background(), id)
	if err != nil {
		return "", err
	}
	return "success", nil
}

// StartTorrentsNow starts one or more torrents immediately disregarding queue position.
func (c *Client) StartTorrentsNow(ctx context.Context, ids ...int) error {
	return c.sendIDsCommand(ctx, "torrent-start-now", ids...)
}

// StopTorrents stops one or more torrents by ID.
func (c *Client) StopTorrents(ctx context.Context, ids ...int) error {
	return c.sendIDsCommand(ctx, "torrent-stop", ids...)
}

// StopTorrent stops a single torrent by ID (backwards-compatible signature).
func (c *Client) StopTorrent(id int) (string, error) {
	err := c.StopTorrents(context.Background(), id)
	if err != nil {
		return "", err
	}
	return "success", nil
}

// VerifyTorrents verifies file checksums for one or more torrents.
func (c *Client) VerifyTorrents(ctx context.Context, ids ...int) error {
	return c.sendIDsCommand(ctx, "torrent-verify", ids...)
}

// VerifyTorrent verifies a single torrent by ID (backwards-compatible signature).
func (c *Client) VerifyTorrent(id int) (string, error) {
	err := c.VerifyTorrents(context.Background(), id)
	if err != nil {
		return "", err
	}
	return "success", nil
}

// ReannounceTorrents forces an immediate tracker re-announce for the given torrents.
func (c *Client) ReannounceTorrents(ctx context.Context, ids ...int) error {
	return c.sendIDsCommand(ctx, "torrent-reannounce", ids...)
}

// StartAll starts all torrents.
func (c *Client) StartAll(ctx ...context.Context) error {
	var reqCtx context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		reqCtx = ctx[0]
	} else {
		reqCtx = context.Background()
	}
	return c.sendIDsCommand(reqCtx, "torrent-start")
}

// StopAll stops all torrents.
func (c *Client) StopAll(ctx ...context.Context) error {
	var reqCtx context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		reqCtx = ctx[0]
	} else {
		reqCtx = context.Background()
	}
	return c.sendIDsCommand(reqCtx, "torrent-stop")
}

// VerifyAll verifies all torrents.
func (c *Client) VerifyAll(ctx ...context.Context) error {
	var reqCtx context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		reqCtx = ctx[0]
	} else {
		reqCtx = context.Background()
	}
	return c.sendIDsCommand(reqCtx, "torrent-verify")
}

type torrentRemoveArgs struct {
	IDs             []int `json:"ids"`
	DeleteLocalData bool  `json:"delete-local-data"`
}

// RemoveTorrents deletes one or more torrents and optionally removes downloaded files.
func (c *Client) RemoveTorrents(ctx context.Context, deleteLocalData bool, ids ...int) error {
	args := torrentRemoveArgs{
		IDs:             ids,
		DeleteLocalData: deleteLocalData,
	}
	return c.Execute(ctx, "torrent-remove", args, nil)
}

// DeleteTorrent deletes a torrent by ID and returns the torrent name (backwards-compatible).
func (c *Client) DeleteTorrent(id int, deleteData bool) (string, error) {
	torrent, err := c.GetTorrent(id)
	name := ""
	if err == nil && torrent != nil {
		name = torrent.Name
	}

	if err := c.RemoveTorrents(context.Background(), deleteData, id); err != nil {
		return "", err
	}
	return name, nil
}

type setLocationArgs struct {
	IDs      []int  `json:"ids"`
	Location string `json:"location"`
	Move     bool   `json:"move"`
}

// SetTorrentLocation sets a new location path for torrents and optionally moves data files.
func (c *Client) SetTorrentLocation(ctx context.Context, location string, move bool, ids ...int) error {
	args := setLocationArgs{
		IDs:      ids,
		Location: location,
		Move:     move,
	}
	return c.Execute(ctx, "torrent-set-location", args, nil)
}

type renamePathArgs struct {
	IDs  []int  `json:"ids"`
	Path string `json:"path"`
	Name string `json:"name"`
}

type RenameResult struct {
	Path string `json:"path"`
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// RenameTorrentPath renames a file or directory within a torrent.
func (c *Client) RenameTorrentPath(ctx context.Context, id int, path, newName string) (*RenameResult, error) {
	args := renamePathArgs{
		IDs:  []int{id},
		Path: path,
		Name: newName,
	}
	var res RenameResult
	if err := c.Execute(ctx, "torrent-rename-path", args, &res); err != nil {
		return nil, fmt.Errorf("failed to rename torrent path: %w", err)
	}
	return &res, nil
}
