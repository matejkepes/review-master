package google_my_business_api

import (
	"bytes"
	"shared-templates"
	"html/template"
)

// CreateReportTemplate returns a template for the monthly report
func CreateReportTemplate() (*template.Template, error) {
	// Create a new template without requiring helper functions
	tmpl := template.New("report")

	// Parse the template
	return tmpl.Parse(shared_templates.MonthlyReportTemplate)
}

// RenderReportToString renders the report template to a string
func RenderReportToString(data *shared_templates.ClientReportData) (string, error) {
	tmpl, err := CreateReportTemplate()
	if err != nil {
		return "", err
	}

	var buf []byte
	buffer := bytes.NewBuffer(buf)
	err = tmpl.Execute(buffer, data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
