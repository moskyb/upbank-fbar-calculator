package upapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

const DefaultHost = "https://api.up.com.au/api/v1"

type Client struct {
	Token      string
	PageSize   int
	Host       string
	HTTPClient *http.Client
	Logger     *slog.Logger
}

type newClientOption func(*Client)

func WithPageSize(pageSize int) newClientOption {
	return func(c *Client) {
		c.PageSize = pageSize
	}
}

func WithHTTPClient(httpClient *http.Client) newClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

func WithLogger(logger *slog.Logger) newClientOption {
	return func(c *Client) {
		c.Logger = logger
	}
}

func WithQuiet() newClientOption {
	return func(c *Client) {
		c.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
}

func WithHost(host string) newClientOption {
	return func(c *Client) {
		c.Host = host
	}
}

func NewClient(token string, opts ...newClientOption) *Client {
	w := os.Stderr
	h := tint.NewHandler(w, &tint.Options{
		NoColor: !isatty.IsTerminal(w.Fd()),
	})
	defaultLogger := slog.New(h)

	c := &Client{
		Token:      token,
		Logger:     defaultLogger,
		HTTPClient: http.DefaultClient,
		Host:       DefaultHost,
		PageSize:   100,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) buildURL(path string) (string, error) {
	path, err := url.JoinPath(c.Host, path)
	if err != nil {
		return "", fmt.Errorf("failed to join path: %w", err)
	}

	return path, nil
}

func (c *Client) makeRequest(req *http.Request) ([]byte, error) {
	q := req.URL.Query()
	q.Add("page[size]", strconv.Itoa(c.PageSize))
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	c.Logger.Info("->", "method", req.Method, "url", req.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	c.Logger.Info("<-", "status", resp.Status, "url", req.URL.String())

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}

		return nil, &errResp
	}

	return body, err
}
