package transmission

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newMockServer(t *testing.T) (*httptest.Server, *int32) {
	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)

		// Verify HTTP Method and Content-Type
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Test CSRF token challenge
		token := r.Header.Get(csrfHeader)
		if token != "valid-csrf-token" {
			w.Header().Set(csrfHeader, "valid-csrf-token")
			w.WriteHeader(http.StatusConflict)
			return
		}

		var req rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		switch req.Method {
		case "session-get":
			res := Session{
				Version:          "4.0.5",
				RPCVersion:       17,
				RPCVersionSemver: "17.0.0",
				DownloadDir:      "/downloads",
				AltSpeedEnabled:  false,
				AltSpeedDown:     50,
				AltSpeedUp:       50,
			}
			data, _ := json.Marshal(res)
			json.NewEncoder(w).Encode(rpcResponse{Result: "success", Arguments: data})

		case "session-stats":
			res := Stats{
				TorrentCount:       10,
				ActiveTorrentCount: 3,
				PausedTorrentCount: 2,
				DownloadSpeed:      1024000,
				UploadSpeed:        512000,
				CurrentStats: currentStats{
					DownloadedBytes: 50000000,
					UploadedBytes:   25000000,
					SecondsActive:   3600,
				},
				CumulativeStats: cumulativeStats{
					DownloadedBytes: 500000000,
					UploadedBytes:   250000000,
					SecondsActive:   86400,
				},
			}
			data, _ := json.Marshal(res)
			json.NewEncoder(w).Encode(rpcResponse{Result: "success", Arguments: data})

		case "torrent-get":
			torrents := Torrents{
				{
					ID:             1,
					Name:           "Ubuntu 22.04 LTS",
					Status:         StatusDownloading,
					PercentDone:    0.75,
					RateDownload:   2048000,
					RateUpload:     102400,
					SizeWhenDone:   3500000000,
					HaveValid:      2625000000,
					UploadRatio:    0.5,
					Eta:            60,
					AddedDate:      1700000000,
					DownloadedEver: 2625000000,
					UploadedEver:   1312500000,
					Trackers: []Tracker{
						{Announce: "http://torrent.ubuntu.com/announce"},
					},
				},
				{
					ID:             2,
					Name:           "Arch Linux 2026",
					Status:         StatusSeeding,
					PercentDone:    1.0,
					RateDownload:   0,
					RateUpload:     512000,
					SizeWhenDone:   1200000000,
					HaveValid:      1200000000,
					UploadRatio:    2.4,
					Eta:            -1,
					AddedDate:      1700001000,
					DownloadedEver: 1200000000,
					UploadedEver:   2880000000,
				},
			}
			res := torrentGetResult{Torrents: torrents}
			data, _ := json.Marshal(res)
			json.NewEncoder(w).Encode(rpcResponse{Result: "success", Arguments: data})

		case "torrent-add":
			res := torrentAddResult{
				TorrentAdded: &TorrentAdded{
					ID:         3,
					Name:       "Debian 12 Netinst",
					HashString: "abcdef1234567890",
				},
			}
			data, _ := json.Marshal(res)
			json.NewEncoder(w).Encode(rpcResponse{Result: "success", Arguments: data})

		case "torrent-start", "torrent-start-now", "torrent-stop", "torrent-verify", "torrent-reannounce", "torrent-remove", "torrent-set-location", "session-set", "session-close", "queue-move-top", "queue-move-up", "queue-move-down", "queue-move-bottom", "torrent-set":
			json.NewEncoder(w).Encode(rpcResponse{Result: "success"})

		case "torrent-rename-path":
			res := RenameResult{Path: "dir/old.txt", Name: "dir/new.txt", ID: 1}
			data, _ := json.Marshal(res)
			json.NewEncoder(w).Encode(rpcResponse{Result: "success", Arguments: data})

		case "free-space":
			res := freeSpaceResult{Path: "/downloads", SizeBytes: 150000000000, TotalSize: 500000000000}
			data, _ := json.Marshal(res)
			json.NewEncoder(w).Encode(rpcResponse{Result: "success", Arguments: data})

		case "port-test":
			res := portTestResult{PortIsOpen: true}
			data, _ := json.Marshal(res)
			json.NewEncoder(w).Encode(rpcResponse{Result: "success", Arguments: data})

		case "blocklist-update":
			res := blocklistResult{BlocklistSize: 250000}
			data, _ := json.Marshal(res)
			json.NewEncoder(w).Encode(rpcResponse{Result: "success", Arguments: data})

		default:
			http.Error(w, "unknown method", http.StatusBadRequest)
		}
	}))

	return server, &requestCount
}

func TestClient_CSRFAutoRenewal(t *testing.T) {
	server, reqCount := newMockServer(t)
	defer server.Close()

	client, err := New(server.URL, "admin", "secret")
	if err != nil {
		t.Fatalf("unexpected error initializing client: %v", err)
	}

	// Should have sent at least 2 requests during initialization (409 challenge + retry)
	if atomic.LoadInt32(reqCount) < 2 {
		t.Errorf("expected at least 2 requests for CSRF negotiation, got %d", atomic.LoadInt32(reqCount))
	}

	if client.Version() != "4.0.5" {
		t.Errorf("expected version 4.0.5, got %q", client.Version())
	}
}

func TestClient_ConcurrentRequests(t *testing.T) {
	server, _ := newMockServer(t)
	defer server.Close()

	client, err := New(server.URL, "admin", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			torrents, err := client.GetTorrents(ctx)
			if err != nil {
				t.Errorf("concurrent GetTorrents failed: %v", err)
			}
			if len(torrents) != 2 {
				t.Errorf("expected 2 torrents, got %d", len(torrents))
			}
		}()
	}
	wg.Wait()
}
