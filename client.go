package transmission

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	defaultRPCURL  = "http://localhost:9091/transmission/rpc"
	csrfHeader     = "X-Transmission-Session-Id"
)

// Option configures a Transmission Client.
type Option func(*Client)

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithTimeout sets the default request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithCustomHeader adds a default HTTP header to all requests.
func WithCustomHeader(key, value string) Option {
	return func(c *Client) {
		c.customHeaders[key] = value
	}
}

// Client interacts with the Transmission RPC API.
type Client struct {
	url           string
	username      string
	password      string
	httpClient    *http.Client
	timeout       time.Duration
	customHeaders map[string]string

	tokenMu sync.RWMutex
	token   string
	sortMu  sync.RWMutex
	sort    Sorting
}

// TransmissionClient is an alias for Client for backwards compatibility.
type TransmissionClient = Client

// New creates a new Transmission RPC Client.
func New(rawURL, username, password string, opts ...Option) (*Client, error) {
	if rawURL == "" {
		rawURL = defaultRPCURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid transmission URL %q: %w", rawURL, err)
	}

	c := &Client{
		url:           parsedURL.String(),
		username:      username,
		password:      password,
		httpClient:    &http.Client{Timeout: defaultTimeout},
		timeout:       defaultTimeout,
		customHeaders: make(map[string]string),
		sort:          SortID,
	}

	for _, opt := range opts {
		opt(c)
	}

	// Test connection and establish CSRF session token
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	if _, err := c.GetSession(ctx); err != nil {
		return c, err
	}

	return c, nil
}

// NewClient is a backwards-compatible constructor.
func NewClient(url, username, password string) *ApiClient {
	c, _ := New(url, username, password)
	return &ApiClient{client: c}
}

// ApiClient is a backwards-compatible wrapper.
type ApiClient struct {
	client *Client
}

func (ac *ApiClient) Post(body string) ([]byte, error) {
	return ac.client.postRaw(context.Background(), []byte(body))
}

// rpcRequest is the JSON-RPC wire request.
type rpcRequest struct {
	Method    string `json:"method"`
	Arguments any    `json:"arguments,omitempty"`
	Tag       int    `json:"tag,omitempty"`
}

// rpcResponse is the JSON-RPC wire response.
type rpcResponse struct {
	Result    string          `json:"result"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
	Tag       int             `json:"tag,omitempty"`
}

// Execute performs a JSON-RPC request to the Transmission daemon with automatic CSRF token renewal.
func (c *Client) Execute(ctx context.Context, method string, args any, result any) error {
	if ctx == nil {
		ctx = context.Background()
	}

	reqBody := rpcRequest{
		Method:    method,
		Arguments: args,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	respBytes, err := c.postRaw(ctx, bodyBytes)
	if err != nil {
		return err
	}

	var resp rpcResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.Result != "success" {
		return errors.New(resp.Result)
	}

	if result != nil && len(resp.Arguments) > 0 {
		if err := json.Unmarshal(resp.Arguments, result); err != nil {
			return fmt.Errorf("failed to decode arguments result: %w", err)
		}
	}

	return nil
}

// postRaw sends a raw POST request, handling HTTP 409 CSRF challenges automatically.
func (c *Client) postRaw(ctx context.Context, body []byte) ([]byte, error) {
	resp, err := c.doRequest(ctx, body)
	if err != nil {
		return nil, err
	}

	// Handle CSRF token challenge
	if resp.StatusCode == http.StatusConflict {
		_ = resp.Body.Close()

		newToken := resp.Header.Get(csrfHeader)
		if newToken != "" {
			c.tokenMu.Lock()
			c.token = newToken
			c.tokenMu.Unlock()
		}

		// Retry once with the new token
		resp, err = c.doRequest(ctx, body)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("transmission RPC error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) doRequest(ctx context.Context, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.customHeaders {
		req.Header.Set(k, v)
	}

	if c.username != "" || c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	c.tokenMu.RLock()
	token := c.token
	c.tokenMu.RUnlock()

	if token != "" {
		req.Header.Set(csrfHeader, token)
	}

	return c.httpClient.Do(req)
}

// Version returns Transmission daemon version string.
func (c *Client) Version() string {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	sess, err := c.GetSession(ctx)
	if err != nil {
		return ""
	}
	return sess.Version
}

// Command is preserved for backwards compatibility.
type Command struct {
	Method    string    `json:"method,omitempty"`
	Arguments arguments `json:"arguments,omitempty"`
	Result    string    `json:"result,omitempty"`
}

type arguments struct {
	Fields             []string        `json:"fields,omitempty"`
	Torrents           Torrents        `json:"torrents,omitempty"`
	Ids                []int           `json:"ids,omitempty"`
	DeleteData         bool            `json:"delete-local-data,omitempty"`
	DownloadDir        string          `json:"download-dir,omitempty"`
	MetaInfo           string          `json:"metainfo,omitempty"`
	Filename           string          `json:"filename,omitempty"`
	TorrentAdded       TorrentAdded    `json:"torrent-added,omitempty"`
	SpeedLimitDown     uint            `json:"speed-limit-down,omitempty"`
	SpeedLimitUp       uint            `json:"speed-limit-up,omitempty"`
	ActiveTorrentCount int             `json:"activeTorrentCount,omitempty"`
	CumulativeStats    cumulativeStats `json:"cumulative-stats,omitempty"`
	CurrentStats       currentStats    `json:"current-stats,omitempty"`
	DownloadSpeed      uint64          `json:"downloadSpeed,omitempty"`
	PausedTorrentCount int             `json:"pausedTorrentCount,omitempty"`
	TorrentCount       int             `json:"torrentCount,omitempty"`
	UploadSpeed        uint64          `json:"uploadSpeed,omitempty"`
	Version            string          `json:"version,omitempty"`
}

// ExecuteCommand executes a legacy Command.
func (c *Client) ExecuteCommand(cmd *Command) (*Command, error) {
	out := &Command{}
	body, err := json.Marshal(cmd)
	if err != nil {
		return out, err
	}
	output, err := c.postRaw(context.Background(), body)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(output, out); err != nil {
		return out, err
	}
	return out, nil
}
