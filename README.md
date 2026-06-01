# growpanel-sdk-go

Official Go SDK for the [GrowPanel](https://growpanel.io) subscription analytics REST API.

```bash
go get github.com/growpanel/growpanel-sdk-go
```

Requires Go 1.22+.

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "os"

    growpanel "github.com/growpanel/growpanel-sdk-go"
)

func main() {
    gp, err := growpanel.New(os.Getenv("GROWPANEL_API_KEY"))
    if err != nil { panic(err) }

    resp, err := gp.WithResponses.GetReportsSummaryWithResponse(context.Background(), nil)
    if err != nil { panic(err) }

    fmt.Println(resp.JSON200.Summary.MrrCurrent)
}
```

## Auth

Get an API key from **app.growpanel.io → Account → API keys**. Pass it to `growpanel.New()`; it's sent as `Authorization: Bearer <key>` on every request.

## Surfaces

`v0.1.0` is a thin wrapper around the generated low-level client. Two flavours are exposed:

| Field | Use when |
|-------|----------|
| `gp.Client.*` | You want the raw `(*http.Response, error)` return — closest to the wire. |
| `gp.WithResponses.*` | You want typed `*<Op>Response` with parsed `JSON200`/`JSON400` fields. **Recommended.** |

Both expose one method per endpoint, named after the operation: `GetReportsMrr`, `PostDataCustomers`, `PutDataCustomersId`, `DeleteIntegrationsWebhooksId`, etc. Browse `internal/generated/growpanel.gen.go` to see the full list, or use the [interactive API reference](https://api.growpanel.io/docs) to discover endpoints.

```go
// Read endpoints
gp.WithResponses.GetReportsSummaryWithResponse(ctx, nil)
gp.WithResponses.GetReportsMrrWithResponse(ctx, &gen.GetReportsMrrParams{Date: ptr("20260101-20260531")})
gp.WithResponses.GetCustomersWithResponse(ctx, &gen.GetCustomersParams{Limit: ptr("50")})

// Write endpoints (POST/PUT/DELETE for /data/*)
gp.WithResponses.PostDataCustomersWithResponse(ctx, gen.PostDataCustomersJSONBody{...})
gp.WithResponses.DeleteDataPlanGroupsIdWithResponse(ctx, "pg_xxx")
```

A small helper makes pointer literals less noisy:

```go
func ptr[T any](v T) *T { return &v }
```

Future versions will add ergonomic group fields (`gp.Reports.GetMrr(...)`, `gp.Data.Customers.Create(...)`) so callers don't need to spell out the full operation ID.

## Configuration

```go
gp, err := growpanel.New(
    os.Getenv("GROWPANEL_API_KEY"),
    growpanel.WithBaseURL("https://api-dev.growpanel.io"),
    growpanel.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
)
```

## Errors

The generated client returns `(*Response, error)`. Network/decoding errors land in `err`; API-level errors (4xx/5xx) come back on the response with status accessible via `resp.StatusCode()`.

```go
resp, err := gp.WithResponses.GetCustomersIdWithResponse(ctx, "cus_doesnotexist", nil)
if err != nil {
    return fmt.Errorf("network: %w", err)
}
switch resp.StatusCode() {
case 200:
    // resp.JSON200 has the typed payload
case 404:
    // customer not found
case 429:
    // rate-limited
}
```

See [error codes](https://growpanel.io/developers/rest-api/error-codes/) for the full list.

## Updating

```bash
go get -u github.com/growpanel/growpanel-sdk-go
```

The SDK is regenerated on every API surface change. Changelog on [GitHub Releases](https://github.com/growpanel/growpanel-sdk-go/releases).

## Related

- [Interactive API reference](https://api.growpanel.io/docs)
- [JavaScript SDK](https://www.npmjs.com/package/@growpanel/sdk)
- [Python SDK](https://pypi.org/project/growpanel/)
- [CLI](https://github.com/growpanel/growpanel-cli)
- [MCP server](https://github.com/growpanel/growpanel-mcp-server)
