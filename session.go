package transmission

import (
	"context"
	"fmt"
	"time"
)

// Session represents global Transmission session settings.
type Session struct {
	AltSpeedDown               int64   `json:"alt-speed-down"`
	AltSpeedEnabled            bool    `json:"alt-speed-enabled"`
	AltSpeedTimeBegin          int     `json:"alt-speed-time-begin"`
	AltSpeedTimeDay            int     `json:"alt-speed-time-day"`
	AltSpeedTimeEnabled        bool    `json:"alt-speed-time-enabled"`
	AltSpeedTimeEnd            int     `json:"alt-speed-time-end"`
	AltSpeedUp                 int64   `json:"alt-speed-up"`
	BlocklistEnabled           bool    `json:"blocklist-enabled"`
	BlocklistSize              int     `json:"blocklist-size"`
	BlocklistURL               string  `json:"blocklist-url"`
	CacheSizeMB                int     `json:"cache-size-mb"`
	ConfigDir                  string  `json:"config-dir"`
	DefaultTrackers            string  `json:"default-trackers"`
	DHTEnabled                 bool    `json:"dht-enabled"`
	DownloadDir                string  `json:"download-dir"`
	DownloadDirFreeSpace       int64   `json:"download-dir-free-space"`
	DownloadQueueEnabled       bool    `json:"download-queue-enabled"`
	DownloadQueueSize          int     `json:"download-queue-size"`
	Encryption                 string  `json:"encryption"`
	IdleSeedingLimit           int     `json:"idle-seeding-limit"`
	IdleSeedingLimitEnabled    bool    `json:"idle-seeding-limit-enabled"`
	IncompleteDir              string  `json:"incomplete-dir"`
	IncompleteDirEnabled       bool    `json:"incomplete-dir-enabled"`
	LPDEnabled                 bool    `json:"lpd-enabled"`
	PeerLimitGlobal            int     `json:"peer-limit-global"`
	PeerLimitPerTorrent        int     `json:"peer-limit-per-torrent"`
	PeerPort                   int     `json:"peer-port"`
	PeerPortRandomOnStart      bool    `json:"peer-port-random-on-start"`
	PexEnabled                 bool    `json:"pex-enabled"`
	PortForwardingEnabled      bool    `json:"port-forwarding-enabled"`
	QueueStalledEnabled        bool    `json:"queue-stalled-enabled"`
	QueueStalledMinutes        int     `json:"queue-stalled-minutes"`
	RenamePartialFiles         bool    `json:"rename-partial-files"`
	RPCVersion                 int     `json:"rpc-version"`
	RPCVersionMinimum          int     `json:"rpc-version-minimum"`
	RPCVersionSemver           string  `json:"rpc-version-semver"`
	ScriptTorrentAddedEnabled  bool    `json:"script-torrent-added-enabled"`
	ScriptTorrentAddedFilename string  `json:"script-torrent-added-filename"`
	ScriptTorrentDoneEnabled   bool    `json:"script-torrent-done-enabled"`
	ScriptTorrentDoneFilename  string  `json:"script-torrent-done-filename"`
	SeedQueueEnabled           bool    `json:"seed-queue-enabled"`
	SeedQueueSize              int     `json:"seed-queue-size"`
	SeedRatioLimit             float64 `json:"seedRatioLimit"`
	SeedRatioLimited           bool    `json:"seedRatioLimited"`
	SessionID                  string  `json:"session-id"`
	SpeedLimitDown             int64   `json:"speed-limit-down"`
	SpeedLimitDownEnabled      bool    `json:"speed-limit-down-enabled"`
	SpeedLimitUp               int64   `json:"speed-limit-up"`
	SpeedLimitUpEnabled        bool    `json:"speed-limit-up-enabled"`
	StartAddedTorrents         bool    `json:"start-added-torrents"`
	TrashOriginalTorrentFiles  bool    `json:"trash-original-torrent-files"`
	UTPEnabled                 bool    `json:"utp-enabled"`
	Version                    string  `json:"version"`
}

// SessionSetOptions defines options for mutating global session settings.
type SessionSetOptions struct {
	AltSpeedDown               *int64   `json:"alt-speed-down,omitempty"`
	AltSpeedEnabled            *bool    `json:"alt-speed-enabled,omitempty"`
	AltSpeedTimeBegin          *int     `json:"alt-speed-time-begin,omitempty"`
	AltSpeedTimeDay            *int     `json:"alt-speed-time-day,omitempty"`
	AltSpeedTimeEnabled        *bool    `json:"alt-speed-time-enabled,omitempty"`
	AltSpeedTimeEnd            *int     `json:"alt-speed-time-end,omitempty"`
	AltSpeedUp                 *int64   `json:"alt-speed-up,omitempty"`
	BlocklistEnabled           *bool    `json:"blocklist-enabled,omitempty"`
	BlocklistURL               *string  `json:"blocklist-url,omitempty"`
	CacheSizeMB                *int     `json:"cache-size-mb,omitempty"`
	DefaultTrackers            *string  `json:"default-trackers,omitempty"`
	DHTEnabled                 *bool    `json:"dht-enabled,omitempty"`
	DownloadDir                *string  `json:"download-dir,omitempty"`
	DownloadQueueEnabled       *bool    `json:"download-queue-enabled,omitempty"`
	DownloadQueueSize          *int     `json:"download-queue-size,omitempty"`
	Encryption                 *string  `json:"encryption,omitempty"`
	IdleSeedingLimit           *int     `json:"idle-seeding-limit,omitempty"`
	IdleSeedingLimitEnabled    *bool    `json:"idle-seeding-limit-enabled,omitempty"`
	IncompleteDir              *string  `json:"incomplete-dir,omitempty"`
	IncompleteDirEnabled       *bool    `json:"incomplete-dir-enabled,omitempty"`
	LPDEnabled                 *bool    `json:"lpd-enabled,omitempty"`
	PeerLimitGlobal            *int     `json:"peer-limit-global,omitempty"`
	PeerLimitPerTorrent        *int     `json:"peer-limit-per-torrent,omitempty"`
	PeerPort                   *int     `json:"peer-port,omitempty"`
	PeerPortRandomOnStart      *bool    `json:"peer-port-random-on-start,omitempty"`
	PexEnabled                 *bool    `json:"pex-enabled,omitempty"`
	PortForwardingEnabled      *bool    `json:"port-forwarding-enabled,omitempty"`
	QueueStalledEnabled        *bool    `json:"queue-stalled-enabled,omitempty"`
	QueueStalledMinutes        *int     `json:"queue-stalled-minutes,omitempty"`
	RenamePartialFiles         *bool    `json:"rename-partial-files,omitempty"`
	ScriptTorrentAddedEnabled  *bool    `json:"script-torrent-added-enabled,omitempty"`
	ScriptTorrentAddedFilename *string  `json:"script-torrent-added-filename,omitempty"`
	ScriptTorrentDoneEnabled   *bool    `json:"script-torrent-done-enabled,omitempty"`
	ScriptTorrentDoneFilename  *string  `json:"script-torrent-done-filename,omitempty"`
	SeedQueueEnabled           *bool    `json:"seed-queue-enabled,omitempty"`
	SeedQueueSize              *int     `json:"seed-queue-size,omitempty"`
	SeedRatioLimit             *float64 `json:"seedRatioLimit,omitempty"`
	SeedRatioLimited           *bool    `json:"seedRatioLimited,omitempty"`
	SpeedLimitDown             *int64   `json:"speed-limit-down,omitempty"`
	SpeedLimitDownEnabled      *bool    `json:"speed-limit-down-enabled,omitempty"`
	SpeedLimitUp               *int64   `json:"speed-limit-up,omitempty"`
	SpeedLimitUpEnabled        *bool    `json:"speed-limit-up-enabled,omitempty"`
	StartAddedTorrents         *bool    `json:"start-added-torrents,omitempty"`
	TrashOriginalTorrentFiles  *bool    `json:"trash-original-torrent-files,omitempty"`
	UTPEnabled                 *bool    `json:"utp-enabled,omitempty"`
}

// GetSession retrieves the current global session settings.
func (c *Client) GetSession(ctx context.Context) (*Session, error) {
	var sess Session
	if err := c.Execute(ctx, "session-get", nil, &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

// SetSession updates global session settings.
func (c *Client) SetSession(ctx context.Context, opt SessionSetOptions) error {
	return c.Execute(ctx, "session-set", opt, nil)
}

// CloseSession cleanly shuts down the Transmission daemon.
func (c *Client) CloseSession(ctx context.Context) error {
	return c.Execute(ctx, "session-close", nil, nil)
}

// Stats holds Transmission cumulative and session stats.
type Stats struct {
	ActiveTorrentCount int             `json:"activeTorrentCount"`
	CumulativeStats    cumulativeStats `json:"cumulative-stats"`
	CurrentStats       currentStats    `json:"current-stats"`
	DownloadSpeed      uint64          `json:"downloadSpeed"`
	PausedTorrentCount int             `json:"pausedTorrentCount"`
	TorrentCount       int             `json:"torrentCount"`
	UploadSpeed        uint64          `json:"uploadSpeed"`
}

type cumulativeStats struct {
	DownloadedBytes uint64        `json:"downloadedBytes"`
	FilesAdded      int           `json:"filesAdded"`
	SecondsActive   time.Duration `json:"secondsActive"`
	SessionCount    int           `json:"sessionCount"`
	UploadedBytes   uint64        `json:"uploadedBytes"`
}

type currentStats struct {
	DownloadedBytes uint64        `json:"downloadedBytes"`
	FilesAdded      int           `json:"filesAdded"`
	SecondsActive   time.Duration `json:"secondsActive"`
	SessionCount    int           `json:"sessionCount"`
	UploadedBytes   uint64        `json:"uploadedBytes"`
}

func (s *Stats) CurrentActiveTime() string {
	return (time.Second * s.CurrentStats.SecondsActive).String()
}

func (s *Stats) CumulativeActiveTime() string {
	return (time.Second * s.CumulativeStats.SecondsActive).String()
}

// GetStats returns current and cumulative Transmission statistics.
func (c *Client) GetStats(ctx ...context.Context) (*Stats, error) {
	var reqCtx context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		reqCtx = ctx[0]
	} else {
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}

	var stats Stats
	if err := c.Execute(reqCtx, "session-stats", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

type freeSpaceArgs struct {
	Path string `json:"path"`
}

type freeSpaceResult struct {
	Path      string `json:"path"`
	SizeBytes int64  `json:"size-bytes"`
	TotalSize int64  `json:"total_size"`
}

// FreeSpace checks remaining and total storage space for the given filesystem path.
func (c *Client) FreeSpace(ctx context.Context, path string) (freeBytes int64, totalBytes int64, err error) {
	var res freeSpaceResult
	if err := c.Execute(ctx, "free-space", freeSpaceArgs{Path: path}, &res); err != nil {
		return 0, 0, fmt.Errorf("free-space check failed: %w", err)
	}
	return res.SizeBytes, res.TotalSize, nil
}

type portTestResult struct {
	PortIsOpen bool `json:"port-is-open"`
}

// PortTest checks if Transmission's incoming peer port is open and reachable from the internet.
func (c *Client) PortTest(ctx context.Context) (bool, error) {
	var res portTestResult
	if err := c.Execute(ctx, "port-test", nil, &res); err != nil {
		return false, fmt.Errorf("port-test failed: %w", err)
	}
	return res.PortIsOpen, nil
}

type blocklistResult struct {
	BlocklistSize int `json:"blocklist-size"`
}

// UpdateBlocklist updates the blocklist from the configured blocklist-url and returns the rule count.
func (c *Client) UpdateBlocklist(ctx context.Context) (int, error) {
	var res blocklistResult
	if err := c.Execute(ctx, "blocklist-update", nil, &res); err != nil {
		return 0, fmt.Errorf("blocklist-update failed: %w", err)
	}
	return res.BlocklistSize, nil
}

// SetAltSpeedEnabled toggles alternative speed limits ("Turtle Mode").
func (c *Client) SetAltSpeedEnabled(ctx context.Context, enabled bool) error {
	return c.SetSession(ctx, SessionSetOptions{
		AltSpeedEnabled: &enabled,
	})
}

// SetSpeedLimit sets the global upload or download speed limit in KB/s.
func (c *Client) SetSpeedLimit(ctx context.Context, limitType SpeedLimitType, limitKB uint) error {
	enabled := true
	limit := int64(limitKB)

	switch limitType {
	case DownloadLimitType:
		return c.SetSession(ctx, SessionSetOptions{
			SpeedLimitDown:        &limit,
			SpeedLimitDownEnabled: &enabled,
		})
	case UploadLimitType:
		return c.SetSession(ctx, SessionSetOptions{
			SpeedLimitUp:        &limit,
			SpeedLimitUpEnabled: &enabled,
		})
	default:
		return fmt.Errorf("unknown speed limit type: %s", limitType)
	}
}

// SetDownloadDir sets the global default download directory.
func (c *Client) SetDownloadDir(ctx context.Context, dir string) error {
	return c.SetSession(ctx, SessionSetOptions{
		DownloadDir: &dir,
	})
}

// NewSessionSetCommand is preserved for backwards compatibility.
func NewSessionSetCommand() *Command {
	return &Command{Method: "session-set"}
}

// SetDownloadDir mutates legacy Command.
func (cmd *Command) SetDownloadDir(dir string) {
	cmd.Arguments.DownloadDir = dir
}

// NewSpeedLimitCommand creates a legacy Command to set speed limits.
func NewSpeedLimitCommand(limitType SpeedLimitType, limit uint) *Command {
	cmd := &Command{Method: "session-set"}
	switch limitType {
	case DownloadLimitType:
		cmd.Arguments.SpeedLimitDown = limit
	case UploadLimitType:
		cmd.Arguments.SpeedLimitUp = limit
	default:
		return nil
	}
	return cmd
}
