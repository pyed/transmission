package transmission

import (
	"sort"
	"strings"
)

// Sorting represents the sort criteria.
type Sorting int

const (
	SortID Sorting = iota
	SortRevID
	SortName
	SortRevName
	SortAge
	SortRevAge
	SortSize
	SortRevSize
	SortProgress
	SortRevProgress
	SortDownSpeed
	SortRevDownSpeed
	SortUpSpeed
	SortRevUpSpeed
	SortDownloaded
	SortRevDownloaded
	SortUploaded
	SortRevUploaded
	SortRatio
	SortRevRatio
)

// SetSort sets the default sorting for GetTorrents calls on the client.
func (c *Client) SetSort(st Sorting) {
	c.sortMu.Lock()
	defer c.sortMu.Unlock()
	c.sort = st
}

func (c *Client) applySort(torrents Torrents, st Sorting) {
	switch st {
	case SortID:
		torrents.SortID(false)
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
	case SortDownSpeed:
		torrents.SortDownSpeed(false)
	case SortRevDownSpeed:
		torrents.SortDownSpeed(true)
	case SortUpSpeed:
		torrents.SortUpSpeed(false)
	case SortRevUpSpeed:
		torrents.SortUpSpeed(true)
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
}

// SortID sorts torrents by ID.
func (t Torrents) SortID(reverse bool) {
	if reverse {
		sort.Slice(t, func(i, j int) bool { return t[i].ID > t[j].ID })
		return
	}
	sort.Slice(t, func(i, j int) bool { return t[i].ID < t[j].ID })
}

// SortName sorts torrents alphabetically by name (case-insensitive).
func (t Torrents) SortName(reverse bool) {
	if reverse {
		sort.Slice(t, func(i, j int) bool { return strings.ToLower(t[i].Name) > strings.ToLower(t[j].Name) })
		return
	}
	sort.Slice(t, func(i, j int) bool { return strings.ToLower(t[i].Name) < strings.ToLower(t[j].Name) })
}

// SortAge sorts torrents by added date.
func (t Torrents) SortAge(reverse bool) {
	if reverse {
		sort.Slice(t, func(i, j int) bool { return t[i].AddedDate > t[j].AddedDate })
		return
	}
	sort.Slice(t, func(i, j int) bool { return t[i].AddedDate < t[j].AddedDate })
}

// SortSize sorts torrents by size when done.
func (t Torrents) SortSize(reverse bool) {
	if reverse {
		sort.Slice(t, func(i, j int) bool { return t[i].SizeWhenDone > t[j].SizeWhenDone })
		return
	}
	sort.Slice(t, func(i, j int) bool { return t[i].SizeWhenDone < t[j].SizeWhenDone })
}

// SortProgress sorts torrents by percent done.
func (t Torrents) SortProgress(reverse bool) {
	if reverse {
		sort.Slice(t, func(i, j int) bool { return t[i].PercentDone > t[j].PercentDone })
		return
	}
	sort.Slice(t, func(i, j int) bool { return t[i].PercentDone < t[j].PercentDone })
}

// SortDownSpeed sorts torrents by download speed.
func (t Torrents) SortDownSpeed(reverse bool) {
	if reverse {
		sort.Slice(t, func(i, j int) bool { return t[i].RateDownload > t[j].RateDownload })
		return
	}
	sort.Slice(t, func(i, j int) bool { return t[i].RateDownload < t[j].RateDownload })
}

// SortUpSpeed sorts torrents by upload speed.
func (t Torrents) SortUpSpeed(reverse bool) {
	if reverse {
		sort.Slice(t, func(i, j int) bool { return t[i].RateUpload > t[j].RateUpload })
		return
	}
	sort.Slice(t, func(i, j int) bool { return t[i].RateUpload < t[j].RateUpload })
}

// SortDownloaded sorts torrents by total downloaded bytes ever.
func (t Torrents) SortDownloaded(reverse bool) {
	if reverse {
		sort.Slice(t, func(i, j int) bool { return t[i].DownloadedEver > t[j].DownloadedEver })
		return
	}
	sort.Slice(t, func(i, j int) bool { return t[i].DownloadedEver < t[j].DownloadedEver })
}

// SortUploaded sorts torrents by total uploaded bytes ever.
func (t Torrents) SortUploaded(reverse bool) {
	if reverse {
		sort.Slice(t, func(i, j int) bool { return t[i].UploadedEver > t[j].UploadedEver })
		return
	}
	sort.Slice(t, func(i, j int) bool { return t[i].UploadedEver < t[j].UploadedEver })
}

// SortRatio sorts torrents by upload ratio.
func (t Torrents) SortRatio(reverse bool) {
	if reverse {
		sort.Slice(t, func(i, j int) bool { return t[i].UploadRatio > t[j].UploadRatio })
		return
	}
	sort.Slice(t, func(i, j int) bool { return t[i].UploadRatio < t[j].UploadRatio })
}
