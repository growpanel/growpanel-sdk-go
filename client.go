// Package growpanel is the official Go SDK for the GrowPanel subscription analytics REST API.
//
// Get started:
//
//	import "github.com/growpanel/growpanel-sdk-go"
//
//	gp, err := growpanel.New("gp_...")
//	if err != nil { ... }
//
//	summary, err := gp.Reports.GetSummary(ctx, nil)
//	fmt.Println(summary.Summary.MrrCurrent)
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

// GrowPanel is the top-level SDK handle. Group fields (Reports, Customers, ...) hold
// the per-area operations; the generated types live under internal/generated and are
// re-exported as needed.
type GrowPanel struct {
	client *gen.Client

	// Group fields are populated in New() and forward to the generated client. Adding a
	// new endpoint to the spec means regenerating then exposing a wrapper here.
	Reports        *reportsGroup
	Customers      *customersGroup
	Plans          *plansGroup
	PlanGroups     *planGroupsGroup
	Segments       *segmentsGroup
	DataSources    *dataSourcesGroup
	DataCustomers  *dataCustomersGroup
	DataInvoices   *dataInvoicesGroup
	DataPlans      *dataPlansGroup
	Profile        *profileGroup
	Notifications  *notificationsGroup
	Webhooks       *webhooksGroup
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

	client, err := gen.NewClient(strings.TrimRight(cfg.BaseURL, "/"),
		gen.WithHTTPClient(cfg.HTTPClient),
		gen.WithRequestEditorFn(authReq))
	if err != nil {
		return nil, fmt.Errorf("growpanel: %w", err)
	}

	gp := &GrowPanel{client: client}
	gp.Reports = &reportsGroup{c: client}
	gp.Customers = &customersGroup{c: client}
	gp.Plans = &plansGroup{c: client}
	gp.PlanGroups = &planGroupsGroup{c: client}
	gp.Segments = &segmentsGroup{c: client}
	gp.DataSources = &dataSourcesGroup{c: client}
	gp.DataCustomers = &dataCustomersGroup{c: client}
	gp.DataInvoices = &dataInvoicesGroup{c: client}
	gp.DataPlans = &dataPlansGroup{c: client}
	gp.Profile = &profileGroup{c: client}
	gp.Notifications = &notificationsGroup{c: client}
	gp.Webhooks = &webhooksGroup{c: client}
	return gp, nil
}

// WithBaseURL overrides the API base URL. Useful for staging environments.
func WithBaseURL(url string) func(*Config) {
	return func(c *Config) { c.BaseURL = url }
}

// WithHTTPClient swaps the underlying HTTP client (custom timeouts, instrumentation, ...).
func WithHTTPClient(hc *http.Client) func(*Config) {
	return func(c *Config) { c.HTTPClient = hc }
}

// Group structs delegate to the generated client. Each group exposes the operations
// belonging to its OpenAPI tag, with stable method names that don\'t change when the
// generator regenerates types.

type reportsGroup       struct{ c *gen.Client }
type customersGroup     struct{ c *gen.Client }
type plansGroup         struct{ c *gen.Client }
type planGroupsGroup    struct{ c *gen.Client }
type segmentsGroup      struct{ c *gen.Client }
type dataSourcesGroup   struct{ c *gen.Client }
type dataCustomersGroup struct{ c *gen.Client }
type dataInvoicesGroup  struct{ c *gen.Client }
type dataPlansGroup     struct{ c *gen.Client }
type profileGroup       struct{ c *gen.Client }
type notificationsGroup struct{ c *gen.Client }
type webhooksGroup      struct{ c *gen.Client }

// Group methods are wired in groups.go (separate file because it\'s mostly
// regeneration-dependent boilerplate). When the spec changes, run `go generate ./...`
// to refresh internal/generated/ and then update groups.go.
