package transmission

import (
	"testing"
)

func sampleTorrents() Torrents {
	return Torrents{
		{ID: 1, Name: "Zebra", AddedDate: 100, SizeWhenDone: 300, PercentDone: 0.2, RateDownload: 1000, RateUpload: 100, DownloadedEver: 50, UploadedEver: 10, UploadRatio: 0.2},
		{ID: 2, Name: "Apple", AddedDate: 200, SizeWhenDone: 100, PercentDone: 0.8, RateDownload: 2000, RateUpload: 500, DownloadedEver: 150, UploadedEver: 300, UploadRatio: 2.0},
		{ID: 3, Name: "Mango", AddedDate: 150, SizeWhenDone: 200, PercentDone: 0.5, RateDownload: 500, RateUpload: 200, DownloadedEver: 100, UploadedEver: 50, UploadRatio: 0.5},
	}
}

func TestTorrents_Sorting(t *testing.T) {
	t.Run("SortID", func(t *testing.T) {
		torrents := sampleTorrents()
		torrents.SortID(false)
		if torrents[0].ID != 1 || torrents[2].ID != 3 {
			t.Errorf("unexpected SortID order: %v", torrents.GetIDs())
		}
		torrents.SortID(true)
		if torrents[0].ID != 3 || torrents[2].ID != 1 {
			t.Errorf("unexpected SortRevID order: %v", torrents.GetIDs())
		}
	})

	t.Run("SortName", func(t *testing.T) {
		torrents := sampleTorrents()
		torrents.SortName(false)
		if torrents[0].Name != "Apple" || torrents[2].Name != "Zebra" {
			t.Errorf("unexpected SortName order: %s, %s, %s", torrents[0].Name, torrents[1].Name, torrents[2].Name)
		}
		torrents.SortName(true)
		if torrents[0].Name != "Zebra" || torrents[2].Name != "Apple" {
			t.Errorf("unexpected SortRevName order: %s, %s, %s", torrents[0].Name, torrents[1].Name, torrents[2].Name)
		}
	})

	t.Run("SortAge", func(t *testing.T) {
		torrents := sampleTorrents()
		torrents.SortAge(false) // oldest (100) first
		if torrents[0].ID != 1 || torrents[2].ID != 2 {
			t.Errorf("unexpected SortAge order: %v", torrents.GetIDs())
		}
		torrents.SortAge(true) // newest (200) first
		if torrents[0].ID != 2 || torrents[2].ID != 1 {
			t.Errorf("unexpected SortRevAge order: %v", torrents.GetIDs())
		}
	})

	t.Run("SortSize", func(t *testing.T) {
		torrents := sampleTorrents()
		torrents.SortSize(false) // smallest (100) first
		if torrents[0].SizeWhenDone != 100 || torrents[2].SizeWhenDone != 300 {
			t.Errorf("unexpected SortSize order")
		}
		torrents.SortSize(true) // largest (300) first
		if torrents[0].SizeWhenDone != 300 || torrents[2].SizeWhenDone != 100 {
			t.Errorf("unexpected SortRevSize order")
		}
	})

	t.Run("SortProgress", func(t *testing.T) {
		torrents := sampleTorrents()
		torrents.SortProgress(false) // lowest progress (0.2) first
		if torrents[0].PercentDone != 0.2 || torrents[2].PercentDone != 0.8 {
			t.Errorf("unexpected SortProgress order")
		}
		torrents.SortProgress(true) // highest progress (0.8) first
		if torrents[0].PercentDone != 0.8 || torrents[2].PercentDone != 0.2 {
			t.Errorf("unexpected SortRevProgress order")
		}
	})

	t.Run("SortDownSpeed", func(t *testing.T) {
		torrents := sampleTorrents()
		torrents.SortDownSpeed(false) // slowest (500) first
		if torrents[0].RateDownload != 500 || torrents[2].RateDownload != 2000 {
			t.Errorf("unexpected SortDownSpeed order")
		}
		torrents.SortDownSpeed(true) // fastest (2000) first
		if torrents[0].RateDownload != 2000 || torrents[2].RateDownload != 500 {
			t.Errorf("unexpected SortRevDownSpeed order")
		}
	})

	t.Run("SortUpSpeed", func(t *testing.T) {
		torrents := sampleTorrents()
		torrents.SortUpSpeed(false) // slowest (100) first
		if torrents[0].RateUpload != 100 || torrents[2].RateUpload != 500 {
			t.Errorf("unexpected SortUpSpeed order")
		}
		torrents.SortUpSpeed(true) // fastest (500) first
		if torrents[0].RateUpload != 500 || torrents[2].RateUpload != 100 {
			t.Errorf("unexpected SortRevUpSpeed order")
		}
	})

	t.Run("SortDownloaded", func(t *testing.T) {
		torrents := sampleTorrents()
		torrents.SortDownloaded(false) // lowest (50) first
		if torrents[0].DownloadedEver != 50 || torrents[2].DownloadedEver != 150 {
			t.Errorf("unexpected SortDownloaded order")
		}
		torrents.SortDownloaded(true) // highest (150) first
		if torrents[0].DownloadedEver != 150 || torrents[2].DownloadedEver != 50 {
			t.Errorf("unexpected SortRevDownloaded order")
		}
	})

	t.Run("SortUploaded", func(t *testing.T) {
		torrents := sampleTorrents()
		torrents.SortUploaded(false) // lowest (10) first
		if torrents[0].UploadedEver != 10 || torrents[2].UploadedEver != 300 {
			t.Errorf("unexpected SortUploaded order")
		}
		torrents.SortUploaded(true) // highest (300) first
		if torrents[0].UploadedEver != 300 || torrents[2].UploadedEver != 10 {
			t.Errorf("unexpected SortRevUploaded order")
		}
	})

	t.Run("SortRatio", func(t *testing.T) {
		torrents := sampleTorrents()
		torrents.SortRatio(false) // lowest (0.2) first
		if torrents[0].UploadRatio != 0.2 || torrents[2].UploadRatio != 2.0 {
			t.Errorf("unexpected SortRatio order")
		}
		torrents.SortRatio(true) // highest (2.0) first
		if torrents[0].UploadRatio != 2.0 || torrents[2].UploadRatio != 0.2 {
			t.Errorf("unexpected SortRevRatio order")
		}
	})
}
