package transmission

// Torrent Status constants (matches libtransmission tr_stat status)
const (
	StatusStopped        = 0 // Torrent is stopped
	StatusCheckPending   = 1 // Queued to check files
	StatusChecking       = 2 // Checking files
	StatusDownloadPending = 3 // Queued to download
	StatusDownloading    = 4 // Downloading
	StatusSeedPending    = 5 // Queued to seed
	StatusSeeding        = 6 // Seeding
)

// Bandwidth Priority constants
const (
	PriorityLow    = -1
	PriorityNormal = 0
	PriorityHigh   = 1
)

// SeedRatioMode constants
const (
	SeedRatioModeGlobal    = 0 // Follow global ratio limit
	SeedRatioModeCustom    = 1 // Follow per-torrent ratio limit
	SeedRatioModeUnlimited = 2 // Unlimited seed ratio
)

// SeedIdleMode constants
const (
	SeedIdleModeGlobal    = 0 // Follow global idle seeding limit
	SeedIdleModeCustom    = 1 // Follow per-torrent idle seeding limit
	SeedIdleModeUnlimited = 2 // Unlimited idle seeding
)

// Error codes reported by Transmission
const (
	ErrorNone           = 0 // No error
	ErrorTrackerWarning = 1 // Tracker warning
	ErrorTrackerError   = 2 // Tracker error
	ErrorLocalError     = 3 // Local filesystem or permission error
)

// SpeedLimitType defines whether a limit applies to download or upload
type SpeedLimitType string

const (
	DownloadLimitType SpeedLimitType = "downloadlimit"
	UploadLimitType   SpeedLimitType = "uploadlimit"
)
