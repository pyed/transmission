package transmission

import (
	"context"
	"testing"
)

func TestSessionAndStats(t *testing.T) {
	server, _ := newMockServer(t)
	defer server.Close()

	client, err := New(server.URL, "admin", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// GetSession
	sess, err := client.GetSession(ctx)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	if sess.Version != "4.0.5" || sess.DownloadDir != "/downloads" {
		t.Errorf("unexpected session: %+v", sess)
	}

	// SetSession
	downDir := "/new/downloads"
	if err := client.SetSession(ctx, SessionSetOptions{DownloadDir: &downDir}); err != nil {
		t.Errorf("SetSession failed: %v", err)
	}

	// SetAltSpeedEnabled
	if err := client.SetAltSpeedEnabled(ctx, true); err != nil {
		t.Errorf("SetAltSpeedEnabled failed: %v", err)
	}

	// SetSpeedLimit
	if err := client.SetSpeedLimit(ctx, DownloadLimitType, 500); err != nil {
		t.Errorf("SetSpeedLimit download failed: %v", err)
	}
	if err := client.SetSpeedLimit(ctx, UploadLimitType, 200); err != nil {
		t.Errorf("SetSpeedLimit upload failed: %v", err)
	}

	// SetDownloadDir
	if err := client.SetDownloadDir(ctx, "/var/downloads"); err != nil {
		t.Errorf("SetDownloadDir failed: %v", err)
	}

	// GetStats
	stats, err := client.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if stats.TorrentCount != 10 || stats.ActiveTorrentCount != 3 {
		t.Errorf("unexpected stats: %+v", stats)
	}
	if stats.CurrentActiveTime() != "1h0m0s" {
		t.Errorf("expected 1h0m0s, got %s", stats.CurrentActiveTime())
	}
	if stats.CumulativeActiveTime() != "24h0m0s" {
		t.Errorf("expected 24h0m0s, got %s", stats.CumulativeActiveTime())
	}

	// FreeSpace
	freeBytes, totalBytes, err := client.FreeSpace(ctx, "/downloads")
	if err != nil {
		t.Fatalf("FreeSpace failed: %v", err)
	}
	if freeBytes != 150000000000 || totalBytes != 500000000000 {
		t.Errorf("unexpected FreeSpace: free=%d, total=%d", freeBytes, totalBytes)
	}

	// PortTest
	open, err := client.PortTest(ctx)
	if err != nil {
		t.Fatalf("PortTest failed: %v", err)
	}
	if !open {
		t.Errorf("expected port open = true")
	}

	// Blocklist
	rules, err := client.UpdateBlocklist(ctx)
	if err != nil {
		t.Fatalf("UpdateBlocklist failed: %v", err)
	}
	if rules != 250000 {
		t.Errorf("expected 250000 rules, got %d", rules)
	}
}
