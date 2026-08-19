package transmission

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"
)

// tracker represents tracker information for a torrent.
type Tracker struct {
	Announce string `json:"announce"`
	ID       int    `json:"id"`
	Scrape   string `json:"scrape"`
	Tier     int    `json:"tier"`
}

type tracker = Tracker // alias for backwards compatibility

// File represents a file inside a torrent.
type File struct {
	BytesCompleted int64  `json:"bytesCompleted"`
	Length         int64  `json:"length"`
	Name           string `json:"name"`
}

// FileStat represents individual file download status.
type FileStat struct {
	BytesCompleted int64 `json:"bytesCompleted"`
	Wanted         bool  `json:"wanted"`
	Priority       int   `json:"priority"`
}

// Peer represents connected peer information.
type Peer struct {
	Address            string  `json:"address"`
	ClientName         string  `json:"clientName"`
	ClientIsChoked     bool    `json:"clientIsChoked"`
	ClientIsInterested bool    `json:"clientIsInterested"`
	FlagStr            string  `json:"flagStr"`
	IsDownloadingFrom  bool    `json:"isDownloadingFrom"`
	IsEncrypted        bool    `json:"isEncrypted"`
	IsIncoming         bool    `json:"isIncoming"`
	IsUploadingTo      bool    `json:"isUploadingTo"`
	PeerIsChoked       bool    `json:"peerIsChoked"`
	PeerIsInterested   bool    `json:"peerIsInterested"`
	Port               int     `json:"port"`
	Progress           float64 `json:"progress"`
	RateToClient       uint64  `json:"rateToClient"`
	RateToPeer         uint64  `json:"rateToPeer"`
}

// Torrent represents all attributes of a Transmission torrent.
type Torrent struct {
	ActivityDate            int64         `json:"activityDate"`
	AddedDate               int64         `json:"addedDate"`
	BandwidthPriority       int           `json:"bandwidthPriority"`
	Comment                 string        `json:"comment"`
	CorruptEver             uint64        `json:"corruptEver"`
	Creator                 string        `json:"creator"`
	DateCreated             int64         `json:"dateCreated"`
	DesiredAvailable        uint64        `json:"desiredAvailable"`
	DoneDate                int64         `json:"doneDate"`
	DownloadDir             string        `json:"downloadDir"`
	DownloadedEver          uint64        `json:"downloadedEver"`
	DownloadLimit           int           `json:"downloadLimit"`
	DownloadLimited         bool          `json:"downloadLimited"`
	EditDate                int64         `json:"editDate"`
	Error                   int           `json:"error"`
	ErrorString             string        `json:"errorString"`
	Eta                     time.Duration `json:"eta"`
	EtaIdle                 time.Duration `json:"etaIdle"`
	FileCount               int           `json:"fileCount"`
	Files                   []File        `json:"files"`
	FileStats               []FileStat    `json:"fileStats"`
	HashString              string        `json:"hashString"`
	HaveUnchecked           uint64        `json:"haveUnchecked"`
	HaveValid               uint64        `json:"haveValid"`
	HonorsSessionLimits     bool          `json:"honorsSessionLimits"`
	ID                      int           `json:"id"`
	IsFinished              bool          `json:"isFinished"`
	IsPrivate               bool          `json:"isPrivate"`
	IsStalled               bool          `json:"isStalled"`
	Labels                  []string      `json:"labels"`
	LeftUntilDone           uint64        `json:"leftUntilDone"`
	ManualAnnounceTime      int64         `json:"manualAnnounceTime"`
	MaxConnectedPeers       int           `json:"maxConnectedPeers"`
	MetadataPercentComplete float64       `json:"metadataPercentComplete"`
	Name                    string        `json:"name"`
	PeerLimit               int           `json:"peer-limit"`
	Peers                   []Peer        `json:"peers"`
	PeersConnected          int           `json:"peersConnected"`
	PeersFrom               map[string]int `json:"peersFrom"`
	PeersGettingFromUs      int           `json:"peersGettingFromUs"`
	PeersSendingToUs        int           `json:"peersSendingToUs"`
	PercentComplete         float64       `json:"percentComplete"`
	PercentDone             float64       `json:"percentDone"`
	PieceCount              int           `json:"pieceCount"`
	PieceSize               int64         `json:"pieceSize"`
	Priorities              []int         `json:"priorities"`
	QueuePosition           int           `json:"queuePosition"`
	RateDownload            uint64        `json:"rateDownload"`
	RateUpload              uint64        `json:"rateUpload"`
	RecheckProgress         float64       `json:"recheckProgress"`
	SecondsDownloading      int           `json:"secondsDownloading"`
	SecondsSeeding          int           `json:"secondsSeeding"`
	SeedIdleLimit           int           `json:"seedIdleLimit"`
	SeedIdleMode            int           `json:"seedIdleMode"`
	SeedRatioLimit          float64       `json:"seedRatioLimit"`
	SeedRatioMode           int           `json:"seedRatioMode"`
	SizeWhenDone            uint64        `json:"sizeWhenDone"`
	StartDate               int64         `json:"startDate"`
	Status                  int           `json:"status"`
	Trackers                []Tracker     `json:"trackers"`
	TotalSize               uint64        `json:"totalSize"`
	TorrentFile             string        `json:"torrentFile"`
	UploadedEver            uint64        `json:"uploadedEver"`
	UploadLimit             int           `json:"uploadLimit"`
	UploadLimited           bool          `json:"uploadLimited"`
	UploadRatio             float64       `json:"uploadRatio"`
	Wanted                  []int         `json:"wanted"`
	Webseeds                []string      `json:"webseeds"`
	WebseedsSendingToUs     int           `json:"webseedsSendingToUs"`
}

// TorrentStatus returns human-readable status name.
func (t *Torrent) TorrentStatus() string {
	switch t.Status {
	case StatusStopped:
		return "Stopped"
	case StatusCheckPending:
		return "Check waiting"
	case StatusChecking:
		return "Checking"
	case StatusDownloadPending:
		return "Download waiting"
	case StatusDownloading:
		return "Downloading"
	case StatusSeedPending:
		return "Seed waiting"
	case StatusSeeding:
		return "Seeding"
	default:
		return "unknown"
	}
}

// Ratio returns formatted upload ratio.
func (t *Torrent) Ratio() string {
	if t.UploadRatio < 0 {
		return "∞"
	}
	return fmt.Sprintf("%.3f", t.UploadRatio)
}

// ETA returns human-readable download ETA.
func (t *Torrent) ETA() string {
	if t.Eta < 0 {
		return "∞"
	}
	return (time.Second * t.Eta).String()
}

// GetTrackers returns announce URLs separated by newlines.
func (t *Torrent) GetTrackers() string {
	var buf bytes.Buffer
	for i := range t.Trackers {
		buf.WriteString(t.Trackers[i].Announce)
		buf.WriteString("\n")
	}
	return buf.String()
}

// Have returns total valid and unchecked downloaded bytes.
func (t *Torrent) Have() uint64 {
	return t.HaveValid + t.HaveUnchecked
}

// Torrents represents a slice of *Torrent with utility methods.
type Torrents []*Torrent

// GetIDs returns a slice of all torrent IDs.
func (t Torrents) GetIDs() []int {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ids
}

var defaultTorrentFields = []string{
	"id", "name", "status", "addedDate", "leftUntilDone", "sizeWhenDone", "totalSize",
	"eta", "uploadRatio", "uploadedEver", "rateDownload", "rateUpload", "downloadDir",
	"hashString", "haveValid", "haveUnchecked", "isFinished", "downloadedEver",
	"percentDone", "seedRatioMode", "seedRatioLimit", "error", "errorString", "trackers",
	"queuePosition", "labels", "bandwidthPriority",
}

// GetTorrents returns all torrents applying the currently configured sort order.
func (c *Client) GetTorrents(ctx ...context.Context) (Torrents, error) {
	var reqCtx context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		reqCtx = ctx[0]
	} else {
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}

	torrents, err := c.GetTorrentsByIDs(reqCtx, nil, defaultTorrentFields...)
	if err != nil {
		return nil, err
	}

	c.sortMu.RLock()
	st := c.sort
	c.sortMu.RUnlock()

	c.applySort(torrents, st)
	return torrents, nil
}

// GetTorrent returns a single torrent by ID.
func (c *Client) GetTorrent(id int, ctx ...context.Context) (*Torrent, error) {
	var reqCtx context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		reqCtx = ctx[0]
	} else {
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}

	torrents, err := c.GetTorrentsByIDs(reqCtx, []int{id}, defaultTorrentFields...)
	if err != nil {
		return nil, err
	}
	if len(torrents) == 0 {
		return nil, errors.New("no torrent with that id")
	}
	return torrents[0], nil
}

type torrentGetArgs struct {
	Fields []string `json:"fields"`
	IDs    any      `json:"ids,omitempty"`
}

type torrentGetResult struct {
	Torrents Torrents `json:"torrents"`
	Removed  []int    `json:"removed,omitempty"`
}

// GetTorrentsByIDs fetches specified torrents with custom field selection.
func (c *Client) GetTorrentsByIDs(ctx context.Context, ids []int, fields ...string) (Torrents, error) {
	if len(fields) == 0 {
		fields = defaultTorrentFields
	}

	args := torrentGetArgs{Fields: fields}
	if len(ids) > 0 {
		args.IDs = ids
	}

	var res torrentGetResult
	if err := c.Execute(ctx, "torrent-get", args, &res); err != nil {
		return nil, err
	}
	return res.Torrents, nil
}

// GetRecentlyActive fetches only torrents that were recently active, plus IDs of recently removed torrents.
func (c *Client) GetRecentlyActive(ctx context.Context, fields ...string) (active Torrents, removedIDs []int, err error) {
	if len(fields) == 0 {
		fields = defaultTorrentFields
	}

	args := torrentGetArgs{
		Fields: fields,
		IDs:    "recently-active",
	}

	var res torrentGetResult
	if err := c.Execute(ctx, "torrent-get", args, &res); err != nil {
		return nil, nil, err
	}
	return res.Torrents, res.Removed, nil
}

// NewGetTorrentsCmd is preserved for backwards compatibility.
func NewGetTorrentsCmd() *Command {
	cmd := &Command{Method: "torrent-get"}
	cmd.Arguments.Fields = defaultTorrentFields
	return cmd
}
