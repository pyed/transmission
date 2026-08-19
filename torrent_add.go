package transmission

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
)

// AddTorrentOptions defines options for adding a new torrent via torrent-add.
type AddTorrentOptions struct {
	Cookies           *string  `json:"cookies,omitempty"`
	DownloadDir       *string  `json:"download-dir,omitempty"`
	Filename          *string  `json:"filename,omitempty"`
	MetaInfo          *string  `json:"metainfo,omitempty"`
	Paused            *bool    `json:"paused,omitempty"`
	PeerLimit         *int     `json:"peer-limit,omitempty"`
	BandwidthPriority *int     `json:"bandwidthPriority,omitempty"`
	FilesUnwanted     []int    `json:"files-unwanted,omitempty"`
	FilesWanted       []int    `json:"files-wanted,omitempty"`
	PriorityHigh      []int    `json:"priority-high,omitempty"`
	PriorityLow       []int    `json:"priority-low,omitempty"`
	PriorityNormal    []int    `json:"priority-normal,omitempty"`
	Labels            []string `json:"labels,omitempty"`
}

// TorrentAdded represents the response from adding a torrent.
type TorrentAdded struct {
	HashString string `json:"hashString"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
}

type torrentAddResult struct {
	TorrentAdded     *TorrentAdded `json:"torrent-added,omitempty"`
	TorrentDuplicate *TorrentAdded `json:"torrent-duplicate,omitempty"`
}

// AddTorrent adds a torrent using custom AddTorrentOptions.
func (c *Client) AddTorrent(ctx context.Context, opt AddTorrentOptions) (*TorrentAdded, error) {
	var res torrentAddResult
	if err := c.Execute(ctx, "torrent-add", opt, &res); err != nil {
		return nil, err
	}

	if res.TorrentAdded != nil {
		return res.TorrentAdded, nil
	}
	if res.TorrentDuplicate != nil {
		return res.TorrentDuplicate, nil
	}

	return nil, fmt.Errorf("torrent-add succeeded but returned no torrent data")
}

// AddTorrentByURL adds a torrent from a magnet link or HTTP/HTTPS URL.
func (c *Client) AddTorrentByURL(ctx context.Context, url string) (*TorrentAdded, error) {
	return c.AddTorrent(ctx, AddTorrentOptions{Filename: &url})
}

// AddTorrentByFile adds a torrent by reading a local .torrent file.
func (c *Client) AddTorrentByFile(ctx context.Context, filePath string) (*TorrentAdded, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read torrent file: %w", err)
	}
	return c.AddTorrentByData(ctx, data)
}

// AddTorrentByData adds a torrent from raw .torrent bytes.
func (c *Client) AddTorrentByData(ctx context.Context, data []byte) (*TorrentAdded, error) {
	b64 := base64.StdEncoding.EncodeToString(data)
	return c.AddTorrent(ctx, AddTorrentOptions{MetaInfo: &b64})
}

// ExecuteAddCommand is preserved for backwards compatibility.
func (c *Client) ExecuteAddCommand(addCmd *Command) (TorrentAdded, error) {
	outCmd, err := c.ExecuteCommand(addCmd)
	if err != nil {
		return TorrentAdded{}, err
	}
	return outCmd.Arguments.TorrentAdded, nil
}

// NewAddCmd creates a legacy add command.
func NewAddCmd() *Command {
	return &Command{Method: "torrent-add"}
}

// NewAddCmdByURL creates a legacy add command from URL or magnet.
func NewAddCmdByURL(url string) *Command {
	cmd := NewAddCmd()
	cmd.Arguments.Filename = url
	return cmd
}

// NewAddCmdByFilename creates a legacy add command with filename.
func NewAddCmdByFilename(filename string) *Command {
	cmd := NewAddCmd()
	cmd.Arguments.Filename = filename
	return cmd
}

// NewAddCmdByFile creates a legacy add command by loading a file.
func NewAddCmdByFile(file string) (*Command, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cmd := NewAddCmd()
	cmd.Arguments.MetaInfo = base64.StdEncoding.EncodeToString(data)
	return cmd, nil
}

func encodeFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
