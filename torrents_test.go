package transmission

import (
	"context"
	"strings"
	"testing"
)

func TestTorrents_GetAndAccessors(t *testing.T) {
	server, _ := newMockServer(t)
	defer server.Close()

	client, err := New(server.URL, "admin", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// GetTorrents
	torrents, err := client.GetTorrents(ctx)
	if err != nil {
		t.Fatalf("GetTorrents failed: %v", err)
	}
	if len(torrents) != 2 {
		t.Fatalf("expected 2 torrents, got %d", len(torrents))
	}

	// GetIDs
	ids := torrents.GetIDs()
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Errorf("unexpected IDs: %v", ids)
	}

	// GetTorrent single
	torrent, err := client.GetTorrent(1, ctx)
	if err != nil {
		t.Fatalf("GetTorrent(1) failed: %v", err)
	}
	if torrent.Name != "Ubuntu 22.04 LTS" {
		t.Errorf("expected Ubuntu, got %s", torrent.Name)
	}

	// Torrent helper methods
	if torrent.TorrentStatus() != "Downloading" {
		t.Errorf("expected 'Downloading', got %q", torrent.TorrentStatus())
	}
	if torrent.Ratio() != "0.500" {
		t.Errorf("expected '0.500', got %q", torrent.Ratio())
	}
	if torrent.ETA() != "1m0s" {
		t.Errorf("expected '1m0s', got %q", torrent.ETA())
	}
	if torrent.Have() != 2625000000 {
		t.Errorf("expected Have() = 2625000000, got %d", torrent.Have())
	}
	if !strings.Contains(torrent.GetTrackers(), "torrent.ubuntu.com") {
		t.Errorf("expected tracker in GetTrackers(), got %s", torrent.GetTrackers())
	}
}

func TestTorrentActions(t *testing.T) {
	server, _ := newMockServer(t)
	defer server.Close()

	client, err := New(server.URL, "admin", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	if err := client.StartTorrents(ctx, 1, 2); err != nil {
		t.Errorf("StartTorrents failed: %v", err)
	}
	if err := client.StartTorrentsNow(ctx, 1); err != nil {
		t.Errorf("StartTorrentsNow failed: %v", err)
	}
	if err := client.StopTorrents(ctx, 1); err != nil {
		t.Errorf("StopTorrents failed: %v", err)
	}
	if err := client.VerifyTorrents(ctx, 1); err != nil {
		t.Errorf("VerifyTorrents failed: %v", err)
	}
	if err := client.ReannounceTorrents(ctx, 1); err != nil {
		t.Errorf("ReannounceTorrents failed: %v", err)
	}
	if err := client.StartAll(ctx); err != nil {
		t.Errorf("StartAll failed: %v", err)
	}
	if err := client.StopAll(ctx); err != nil {
		t.Errorf("StopAll failed: %v", err)
	}
	if err := client.VerifyAll(ctx); err != nil {
		t.Errorf("VerifyAll failed: %v", err)
	}
	if err := client.RemoveTorrents(ctx, true, 1); err != nil {
		t.Errorf("RemoveTorrents failed: %v", err)
	}
	if err := client.SetTorrentLocation(ctx, "/new/path", true, 1); err != nil {
		t.Errorf("SetTorrentLocation failed: %v", err)
	}

	res, err := client.RenameTorrentPath(ctx, 1, "dir/old.txt", "dir/new.txt")
	if err != nil {
		t.Errorf("RenameTorrentPath failed: %v", err)
	}
	if res.Name != "dir/new.txt" {
		t.Errorf("expected new.txt, got %s", res.Name)
	}
}

func TestTorrentAdding(t *testing.T) {
	server, _ := newMockServer(t)
	defer server.Close()

	client, err := New(server.URL, "admin", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	added, err := client.AddTorrentByURL(ctx, "magnet:?xt=urn:btih:xyz")
	if err != nil {
		t.Fatalf("AddTorrentByURL failed: %v", err)
	}
	if added.ID != 3 || added.Name != "Debian 12 Netinst" {
		t.Errorf("unexpected added torrent: %+v", added)
	}

	addedData, err := client.AddTorrentByData(ctx, []byte("d8:announce..."))
	if err != nil {
		t.Fatalf("AddTorrentByData failed: %v", err)
	}
	if addedData.ID != 3 {
		t.Errorf("expected ID 3, got %d", addedData.ID)
	}
}

func TestTorrentMutator(t *testing.T) {
	server, _ := newMockServer(t)
	defer server.Close()

	client, err := New(server.URL, "admin", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	downLimit := 500
	limited := true

	err = client.SetTorrent(ctx, []int{1}, TorrentSetOptions{
		DownloadLimit:   &downLimit,
		DownloadLimited: &limited,
		Labels:          []string{"linux", "iso"},
	})
	if err != nil {
		t.Errorf("SetTorrent failed: %v", err)
	}
}
