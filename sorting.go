package transmission

import "sort"

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
	SortDownloaded
	SortRevDownloaded
	SortUploaded
	SortRevUploaded
	SortRatio
	SortRevRatio
)

// sorting types
type (
	byID         Torrents
	byName       Torrents
	byAge        Torrents
	bySize       Torrents
	byProgress   Torrents
	byDownloaded Torrents
	byUploaded   Torrents
	byRatio      Torrents
)

func (t byID) Len() int           { return len(t) }
func (t byID) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byID) Less(i, j int) bool { return t[i].ID < t[j].ID }

func (t byName) Len() int           { return len(t) }
func (t byName) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byName) Less(i, j int) bool { return t[i].Name < t[j].Name }

func (t byAge) Len() int           { return len(t) }
func (t byAge) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byAge) Less(i, j int) bool { return t[i].AddedDate < t[j].AddedDate }

func (t bySize) Len() int           { return len(t) }
func (t bySize) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t bySize) Less(i, j int) bool { return t[i].SizeWhenDone < t[j].SizeWhenDone }

func (t byProgress) Len() int           { return len(t) }
func (t byProgress) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byProgress) Less(i, j int) bool { return t[i].PercentDone < t[j].PercentDone }

func (t byDownloaded) Len() int           { return len(t) }
func (t byDownloaded) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byDownloaded) Less(i, j int) bool { return t[i].DownloadedEver < t[j].DownloadedEver }

func (t byUploaded) Len() int           { return len(t) }
func (t byUploaded) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byUploaded) Less(i, j int) bool { return t[i].UploadedEver < t[j].UploadedEver }

func (t byRatio) Len() int           { return len(t) }
func (t byRatio) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byRatio) Less(i, j int) bool { return t[i].UploadRatio < t[j].UploadRatio }

func (t Torrents) SortID(reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(byID(t)))
		return
	}
	sort.Sort(byID(t))
}

func (t Torrents) SortName(reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(byName(t)))
		return
	}
	sort.Sort(byName(t))
}

func (t Torrents) SortAge(reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(byAge(t)))
		return
	}
	sort.Sort(byAge(t))
}

func (t Torrents) SortSize(reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(bySize(t)))
		return
	}
	sort.Sort(bySize(t))
}

func (t Torrents) SortProgress(reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(byProgress(t)))
		return
	}
	sort.Sort(byProgress(t))
}

func (t Torrents) SortDownloaded(reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(byDownloaded(t)))
		return
	}
	sort.Sort(byDownloaded(t))
}

func (t Torrents) SortUploaded(reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(byUploaded(t)))
		return
	}
	sort.Sort(byUploaded(t))
}

func (t Torrents) SortRatio(reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(byRatio(t)))
		return
	}
	sort.Sort(byRatio(t))
}
