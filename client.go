// Package growpanel is the official Go SDK for the GrowPanel subscription analytics REST API.
//
// v0.1.0 ships a thin wrapper around the generated low-level client. Every operation in the
// OpenAPI spec is reachable as a typed method via gp.Client (raw response) or gp.WithResponses
// (typed *Response wrapper). Future versions will add ergonomic group fields (gp.Reports.*,
// gp.Data.Customers.*, …) — for now consume the generated client directly:
//
//	gp, err := growpanel.New("gp_...")
//	if err != nil { ... }
//
//	resp, err := gp.WithResponses.GetReportsSummaryWithResponse(ctx, nil)
//	if err != nil { ... }
//	fmt.Println(resp.JSON200.Summary.MrrCurrent)
package growpanel

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	gen "github.com/growpanel/growpanel-sdk-go/internal/generated"
)

const DefaultBaseURL = "https://api.growpanel.io"

// Config controls how a GrowPanel client is constructed.
type Config struct {
	APIKey     string       // Required. Issued by /account/api-keys.
	BaseURL    string       // Optional. Defaults to DefaultBaseURL.
	HTTPClient *http.Client // Optional. Defaults to http.DefaultClient.
}

// GrowPanel exposes two flavours of the same client:
//   - Client: raw operations returning (*http.Response, error). Closest to the wire.
//   - WithResponses: typed-response wrappers returning *<Op>Response with parsed JSON fields.
//     This is what you usually want.
type GrowPanel struct {
	Client        *gen.Client
	WithResponses *gen.ClientWithResponses
}

// New builds a GrowPanel client. The API key is sent on every request as
// `Authorization: Bearer <key>`.
func New(apiKey string, opts ...func(*Config)) (*GrowPanel, error) {
	cfg := Config{APIKey: apiKey, BaseURL: DefaultBaseURL, HTTPClient: http.DefaultClient}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("growpanel: APIKey is required")
	}

	authReq := func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		return nil
	}

	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	raw, err := gen.NewClient(baseURL,
		gen.WithHTTPClient(cfg.HTTPClient),
		gen.WithRequestEditorFn(authReq))
	if err != nil {
		return nil, fmt.Errorf("growpanel: %w", err)
	}
	withResp, err := gen.NewClientWithResponses(baseURL,
		gen.WithHTTPClient(cfg.HTTPClient),
		gen.WithRequestEditorFn(authReq))
	if err != nil {
		return nil, fmt.Errorf("growpanel: %w", err)
	}

	return &GrowPanel{Client: raw, WithResponses: withResp}, nil
}

// WithBaseURL overrides the API base URL. Useful for staging environments.
func WithBaseURL(url string) func(*Config) {
	return func(c *Config) { c.BaseURL = url }
}

// WithHTTPClient swaps the underlying HTTP client (custom timeouts, instrumentation, ...).
func WithHTTPClient(hc *http.Client) func(*Config) {
	return func(c *Config) { c.HTTPClient = hc }
}
