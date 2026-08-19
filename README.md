# transmission

A modern, idiomatic Go client library for the [Transmission BitTorrent client](https://transmissionbt.com/) RPC API.

[![Go Reference](https://pkg.go.dev/badge/github.com/pyed/transmission.svg)](https://pkg.go.dev/github.com/pyed/transmission)

## Features

- **Full Transmission RPC Coverage**: Supports Transmission 2.x, 3.x, and 4.x.
- **Context Support**: Complete `context.Context` propagation on all methods.
- **Thread-safe & Resilient**: Handles CSRF session token (`X-Transmission-Session-Id`) negotiation and HTTP 409 re-authentication automatically with thread-safe read/write locking.
- **Rich Methods**:
  - **Torrents**: Start, StartNow, Stop, Verify, Reannounce, SetLocation, RenamePath, Remove, Add (Magnet/URL/File/Bytes), Set (Speed limits, Priority, Wanted files, Labels, Ratio).
  - **Session**: Session settings, Alternative Speed Limits ("Turtle Mode"), Download Directory, Speed Limits, Daemon Stats, Shutdown.
  - **Utilities**: Free space check (`FreeSpace`), Port testing (`PortTest`), Blocklist update (`UpdateBlocklist`), Queue manipulation (`QueueMoveTop`, `QueueMoveUp`, `QueueMoveDown`, `QueueMoveBottom`).
- **Zero External Dependencies**: Uses only standard library Go.

## Installation

```bash
go get github.com/pyed/transmission
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pyed/transmission"
)

func main() {
	ctx := context.Background()

	// Initialize client
	client, err := transmission.New("http://localhost:9091/transmission/rpc", "username", "password")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	// List torrents
	torrents, err := client.GetTorrents(ctx)
	if err != nil {
		log.Fatalf("failed to get torrents: %v", err)
	}

	for _, t := range torrents {
		fmt.Printf("[%d] %s - %s (%.1f%%)\n", t.ID, t.Name, t.TorrentStatus(), t.PercentDone*100)
	}

	// Add a torrent
	added, err := client.AddTorrentByURL(ctx, "magnet:?xt=urn:btih:...")
	if err != nil {
		log.Fatalf("failed to add torrent: %v", err)
	}
	fmt.Printf("Added torrent: %s (ID: %d)\n", added.Name, added.ID)

	// Toggle Turtle Mode
	if err := client.SetAltSpeedEnabled(ctx, true); err != nil {
		log.Fatalf("failed to enable turtle mode: %v", err)
	}

	// Check free disk space
	free, total, err := client.FreeSpace(ctx, "/downloads")
	if err == nil {
		fmt.Printf("Free disk space: %d GB / %d GB\n", free/(1024*1024*1024), total/(1024*1024*1024))
	}
}
```

## License

MIT
