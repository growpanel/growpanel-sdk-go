# growpanel-sdk-go

Official Go SDK for the [GrowPanel](https://growpanel.io) subscription analytics REST API.

```bash
go get github.com/growpanel/growpanel-sdk-go
```

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

    resp, err := gp.Reports.GetSummary(context.Background(), nil)
    if err != nil { panic(err) }

    fmt.Println(resp.JSON200.Summary.MrrCurrent)
}
```

## Auth

Get an API key from **app.growpanel.io → Account → API keys**. Pass it to `growpanel.New()`; it\'s sent as `Authorization: Bearer <key>`.

## First-time setup (generate from spec)

The repository ships without `internal/generated/` — generate it from the live OpenAPI spec:

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
curl -s https://api-dev.growpanel.io/openapi.json > openapi.json
oapi-codegen -config oapi-codegen.yaml openapi.json
go mod tidy
go build ./...
```

In CI, the [SDK pipeline](../growpanel-api/.github/workflows/sdk-pipeline.yml) does this automatically on every API change.

## Surfaces

Operations are grouped by API area:

- `gp.Reports.*` — MRR, leads, cohorts, cashflow, retention, churn
- `gp.Customers.*` — list + detail (analytics view)
- `gp.Plans.*` — list plans
- `gp.PlanGroups.*`, `gp.Segments.*`, `gp.DataSources.*`, `gp.DataCustomers.*`, `gp.DataInvoices.*`, `gp.DataPlans.*` — data management with full CRUD
- `gp.Profile.*`, `gp.Notifications.*`, `gp.Webhooks.*` — account & integrations

For anything not exposed as a curated method, the underlying generated client is available — see `internal/generated/`.

## Interactive docs

The full reference (with try-it-out) lives at [api.growpanel.io/docs](https://api.growpanel.io/docs).
