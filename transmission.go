package transmission

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	StatusStopped = iota
	StatusCheckPending
	StatusChecking
	StatusDownloadPending
	StatusDownloading
	StatusSeedPending
	StatusSeeding
)

//TransmissionClient to talk to transmission
type TransmissionClient struct {
	apiclient *ApiClient
}

type Command struct {
	Method    string    `json:"method,omitempty"`
	Arguments arguments `json:"arguments,omitempty"`
	Result    string    `json:"result,omitempty"`
}

type arguments struct {
	Fields       []string     `json:"fields,omitempty"`
	Torrents     Torrents     `json:"torrents,omitempty"`
	Ids          []int        `json:"ids,omitempty"`
	DeleteData   bool         `json:"delete-local-data,omitempty"`
	DownloadDir  string       `json:"download-dir,omitempty"`
	MetaInfo     string       `json:"metainfo,omitempty"`
	Filename     string       `json:"filename,omitempty"`
	TorrentAdded TorrentAdded `json:"torrent-added"`
	// Stats
	ActiveTorrentCount int             `json:"activeTorrentCount"`
	CumulativeStats    cumulativeStats `json:"cumulative-stats"`
	CurrentStats       currentStats    `json:"current-stats"`
	DownloadSpeed      uint64          `json:"downloadSpeed"`
	PausedTorrentCount int             `json:"pausedTorrentCount"`
	TorrentCount       int             `json:"torrentCount"`
	UploadSpeed        uint64          `json:"uploadSpeed"`
	Version            string          `json:"version"`
}

type tracker struct {
	Announce string `json:"announce"`
	Id       int    `json:"id"`
	Scrape   string `json:"scrape"`
	Tire     int    `json:"tire"`
}

//TorrentAdded data returning
type TorrentAdded struct {
	HashString string `json:"hashString"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
}

// session-stats
type Stats struct {
	ActiveTorrentCount int
	CumulativeStats    cumulativeStats
	CurrentStats       currentStats
	DownloadSpeed      uint64
	PausedTorrentCount int
	TorrentCount       int
	UploadSpeed        uint64
}
type cumulativeStats struct {
	DownloadedBytes uint64 `json:"downloadedBytes"`
	FilesAdded      int    `json:"filesAdded"`
	SecondsActive   int    `json:"secondsActive"`
	SessionCount    int    `json:"sessionCount"`
	UploadedBytes   uint64 `json:"uploadedBytes"`
}
type currentStats struct {
	DownloadedBytes uint64 `json:"downloadedBytes"`
	FilesAdded      int    `json:"filesAdded"`
	SecondsActive   int    `json:"secondsActive"`
	SessionCount    int    `json:"sessionCount"`
	UploadedBytes   uint64 `json:"uploadedBytes"`
}

//Torrent struct for torrents
type Torrent struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Status         int       `json:"status"`
	AddedDate      int64     `json:"addedDate"`
	LeftUntilDone  uint64    `json:"leftUntilDone"`
	SizeWhenDone   uint64    `json:"sizeWhenDone"`
	Eta            int       `json:"eta"`
	UploadRatio    float64   `json:"uploadRatio"`
	RateDownload   uint64    `json:"rateDownload"`
	RateUpload     uint64    `json:"rateUpload"`
	DownloadDir    string    `json:"downloadDir"`
	DownloadedEver uint64    `json:"downloadedEver"`
	UploadedEver   uint64    `json:"uploadedEver"`
	IsFinished     bool      `json:"isFinished"`
	PercentDone    float64   `json:"percentDone"`
	SeedRatioMode  int       `json:"seedRatioMode"`
	Trackers       []tracker `json:"trackers"`
	Error          int       `json:"error"`
	ErrorString    string    `json:"errorString"`
}

// Status translates the status of the torrent
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

// Ratio returns the upload ratio of the torrent
func (t *Torrent) Ratio() string {
	if t.UploadRatio < 0 {
		return "∞"
	}
	return fmt.Sprintf("%.3f", t.UploadRatio)
}

// ETA returns the time left for the download to finish
func (t *Torrent) ETA() string {
	if t.Eta < 0 {
		return "∞"
	}
	return fmt.Sprintf("%d", t.Eta)
}

// Torrents represent []Torrent
type Torrents []*Torrent

// GetIDs returns []int of all the ids
func (t Torrents) GetIDs() []int {
	ids := make([]int, 0, len(t))
	for i := range t {
		ids = append(ids, t[i].ID)
	}
	return ids
}

// sortType keeps track of which sorting we are using
var sortType = SortID // SortID is transmission's default

// SetSort takes a 'Sorting' to set 'sortType'
func (ac *TransmissionClient) SetSort(st Sorting) {
	sortType = st
}

//New create new transmission torrent
func New(url string, username string, password string) *TransmissionClient {
	apiclient := NewClient(url, username, password)
	return &TransmissionClient{apiclient: apiclient}
}

//GetTorrents get a list of torrents
func (ac *TransmissionClient) GetTorrents() (Torrents, error) {
	cmd := NewGetTorrentsCmd()

	out, err := ac.ExecuteCommand(cmd)
	if err != nil {
		return nil, err
	}

	torrents := out.Arguments.Torrents

	// sorting
	switch sortType {
	case SortID:
		return torrents, nil // already sorted by ID
	case SortRevID:
		torrents.SortID(true)
	case SortName:
		torrents.SortName(false)
	case SortRevName:
		torrents.SortName(true)
	case SortAge:
		torrents.SortAge(false)
	case SortRevAge:
		torrents.SortAge(true)
	case SortSize:
		torrents.SortSize(false)
	case SortRevSize:
		torrents.SortSize(true)
	case SortProgress:
		torrents.SortProgress(false)
	case SortRevProgress:
		torrents.SortProgress(true)
	case SortDownloaded:
		torrents.SortDownloaded(false)
	case SortRevDownloaded:
		torrents.SortDownloaded(true)
	case SortUploaded:
		torrents.SortUploaded(false)
	case SortRevUploaded:
		torrents.SortUploaded(true)
	case SortRatio:
		torrents.SortRatio(false)
	case SortRevRatio:
		torrents.SortRatio(true)
	}

	return torrents, nil
}

// GetTorrent takes an id and returns *Torrent
func (ac *TransmissionClient) GetTorrent(id int) (*Torrent, error) {
	cmd := NewGetTorrentsCmd()
	cmd.Arguments.Ids = append(cmd.Arguments.Ids, id)

	out, err := ac.ExecuteCommand(cmd)
	if err != nil {
		return &Torrent{}, err
	}

	if len(out.Arguments.Torrents) > 0 {
		return out.Arguments.Torrents[0], nil
	}
	return &Torrent{}, errors.New("No torrent with that id")
}

// Delete takes a bool, if true it will delete with data;
// returns the name of the deleted torrent if it succeed
func (ac *TransmissionClient) DeleteTorrent(id int, wd bool) (string, error) {
	torrent, err := ac.GetTorrent(id)
	if err != nil {
		return "", err
	}

	cmd := newDelCmd(id, wd)

	_, err = ac.ExecuteCommand(cmd)
	if err != nil {
		return "", err
	}

	return torrent.Name, nil
}

// GetStats returns "session-stats"
func (ac *TransmissionClient) GetStats() (*Stats, error) {
	cmd := &Command{
		Method: "session-stats",
	}

	out, err := ac.ExecuteCommand(cmd)
	if err != nil {
		return nil, err
	}

	return &Stats{
		ActiveTorrentCount: out.Arguments.ActiveTorrentCount,
		CumulativeStats:    out.Arguments.CumulativeStats,
		CurrentStats:       out.Arguments.CurrentStats,
		DownloadSpeed:      out.Arguments.DownloadSpeed,
		PausedTorrentCount: out.Arguments.PausedTorrentCount,
		TorrentCount:       out.Arguments.TorrentCount,
		UploadSpeed:        out.Arguments.UploadSpeed,
	}, nil
}

//StartTorrent start the torrent
func (ac *TransmissionClient) StartTorrent(id int) (string, error) {
	return ac.sendSimpleCommand("torrent-start", id)
}

//StopTorrent start the torrent
func (ac *TransmissionClient) StopTorrent(id int) (string, error) {
	return ac.sendSimpleCommand("torrent-stop", id)
}

// VerifyTorrent verifies a torrent
func (ac *TransmissionClient) VerifyTorrent(id int) (string, error) {
	return ac.sendSimpleCommand("torrent-verify", id)
}

// StartAll starts all the torrents
func (ac *TransmissionClient) StartAll() error {
	cmd := Command{Method: "torrent-start"}
	torrents, err := ac.GetTorrents()
	if err != nil {
		return err
	}

	cmd.Arguments.Ids = torrents.GetIDs()
	if _, err := ac.sendCommand(cmd); err != nil {
		return err
	}

	return nil
}

// StopAll stops all torrents
func (ac *TransmissionClient) StopAll() error {
	cmd := Command{Method: "torrent-stop"}
	torrents, err := ac.GetTorrents()
	if err != nil {
		return err
	}

	cmd.Arguments.Ids = torrents.GetIDs()
	if _, err := ac.sendCommand(cmd); err != nil {
		return err
	}

	return nil
}

// VerifyAll verfies all torrents
func (ac *TransmissionClient) VerifyAll() error {
	cmd := Command{Method: "torrent-verify"}

	torrents, err := ac.GetTorrents()
	if err != nil {
		return err
	}

	cmd.Arguments.Ids = torrents.GetIDs()
	if _, err := ac.sendCommand(cmd); err != nil {
		return err
	}

	return nil
}

func NewGetTorrentsCmd() *Command {
	cmd := &Command{}

	cmd.Method = "torrent-get"
	cmd.Arguments.Fields = []string{"id", "name",
		"status", "addedDate", "leftUntilDone", "sizeWhenDone", "eta", "uploadRatio", "uploadedEver",
		"rateDownload", "rateUpload", "downloadDir", "isFinished", "downloadedEver",
		"percentDone", "seedRatioMode", "error", "errorString", "trackers"}

	return cmd
}

func NewAddCmd() *Command {
	cmd := &Command{}
	cmd.Method = "torrent-add"
	return cmd
}

// URL or magnet
func NewAddCmdByURL(url string) *Command {
	cmd := NewAddCmd()
	cmd.Arguments.Filename = url
	return cmd
}

func NewAddCmdByFilename(filename string) *Command {
	cmd := NewAddCmd()
	cmd.Arguments.Filename = filename
	return cmd
}

func NewAddCmdByFile(file string) (*Command, error) {
	cmd := NewAddCmd()

	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cmd.Arguments.MetaInfo = base64.StdEncoding.EncodeToString(fileData)

	return cmd, nil
}

func (cmd *Command) SetDownloadDir(dir string) {
	cmd.Arguments.DownloadDir = dir
}

func newDelCmd(id int, removeFile bool) *Command {
	cmd := &Command{}
	cmd.Method = "torrent-remove"
	cmd.Arguments.Ids = []int{id}
	cmd.Arguments.DeleteData = removeFile
	return cmd
}

func (ac *TransmissionClient) ExecuteCommand(cmd *Command) (*Command, error) {
	out := &Command{}

	body, err := json.Marshal(cmd)
	if err != nil {
		return out, err
	}
	output, err := ac.apiclient.Post(string(body))
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(output, &out)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (ac *TransmissionClient) ExecuteAddCommand(addCmd *Command) (TorrentAdded, error) {
	outCmd, err := ac.ExecuteCommand(addCmd)
	if err != nil {
		return TorrentAdded{}, err
	}
	return outCmd.Arguments.TorrentAdded, nil
}

func encodeFile(file string) (string, error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(fileData), nil
}

// Version returns transmission's version
func (ac *TransmissionClient) Version() string {
	cmd := Command{Method: "session-get"}

	resp, _ := ac.sendCommand(cmd)
	return resp.Arguments.Version
}

func (ac *TransmissionClient) sendSimpleCommand(method string, id int) (result string, err error) {
	cmd := Command{Method: method}
	cmd.Arguments.Ids = []int{id}
	resp, err := ac.sendCommand(cmd)
	return resp.Result, err
}

func (ac *TransmissionClient) sendCommand(cmd Command) (response Command, err error) {
	body, err := json.Marshal(cmd)
	if err != nil {
		return
	}
	output, err := ac.apiclient.Post(string(body))
	if err != nil {
		return
	}
	err = json.Unmarshal(output, &response)
	if err != nil {
		return
	}
	return response, nil
}
