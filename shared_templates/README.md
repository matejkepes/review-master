# Shared Templates Module

This module contains shared HTML templates and data structures used across multiple services in the Review Master system.

## Purpose

Eliminates code duplication by providing a single source of truth for HTML templates used by:
- `google_my_business` service (PDF report generation)
- `rm_client_portal` service (web report viewing)

## Files

- `monthly_report.go` - Contains the `MonthlyReportTemplate` constant with 1200+ line HTML/CSS template
- `types.go` - Shared data structures used by both services
- `template_test.go` - Basic tests to ensure template parses and executes correctly
- `README.md` - This documentation

## Usage

```go
import "shared_templates"

// Parse the template
tmpl, err := template.New("monthly_report").Parse(shared_templates.MonthlyReportTemplate)
if err != nil {
    return err
}

// Execute with your data
var data shared_templates.ClientReportData
// ... populate data ...

err = tmpl.Execute(writer, data)
```

## Integration

Both consuming services include this module via:

```go
// In go.mod
require shared_templates v0.0.0
replace shared_templates => ../shared_templates
```

## Maintenance

When updating the template:
1. Only modify files in this module
2. Run tests: `go test`
3. Both services will automatically use the updated template
4. No need to update multiple locations