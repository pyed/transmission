package transmission

import (
	"context"
)

// TorrentSetOptions specifies parameters for mutating torrents via torrent-set.
type TorrentSetOptions struct {
	BandwidthPriority             *int      `json:"bandwidthPriority,omitempty"`
	DownloadLimit                 *int      `json:"downloadLimit,omitempty"`
	DownloadLimited               *bool     `json:"downloadLimited,omitempty"`
	FilesUnwanted                 []int     `json:"files-unwanted,omitempty"`
	FilesWanted                   []int     `json:"files-wanted,omitempty"`
	Group                         *string   `json:"group,omitempty"`
	HonorsSessionLimits           *bool     `json:"honorsSessionLimits,omitempty"`
	Labels                        []string  `json:"labels,omitempty"`
	Location                      *string   `json:"location,omitempty"`
	PeerLimit                     *int      `json:"peer-limit,omitempty"`
	PriorityHigh                  []int     `json:"priority-high,omitempty"`
	PriorityLow                   []int     `json:"priority-low,omitempty"`
	PriorityNormal                []int     `json:"priority-normal,omitempty"`
	QueuePosition                 *int      `json:"queuePosition,omitempty"`
	SeedIdleLimit                 *int      `json:"seedIdleLimit,omitempty"`
	SeedIdleMode                  *int      `json:"seedIdleMode,omitempty"`
	SeedRatioLimit                *float64  `json:"seedRatioLimit,omitempty"`
	SeedRatioMode                 *int      `json:"seedRatioMode,omitempty"`
	SequentialDownload            *bool     `json:"sequentialDownload,omitempty"`
	SequentialDownloadFromPiece   *int      `json:"sequentialDownloadFromPiece,omitempty"`
	TrackerList                   *string   `json:"trackerList,omitempty"`
	UploadLimit                   *int      `json:"uploadLimit,omitempty"`
	UploadLimited                 *bool     `json:"uploadLimited,omitempty"`
}

type torrentSetArgs struct {
	TorrentSetOptions
	IDs []int `json:"ids,omitempty"`
}

// SetTorrent mutates settings for the specified torrents.
func (c *Client) SetTorrent(ctx context.Context, ids []int, opt TorrentSetOptions) error {
	args := torrentSetArgs{
		TorrentSetOptions: opt,
		IDs:               ids,
	}
	return c.Execute(ctx, "torrent-set", args, nil)
}
