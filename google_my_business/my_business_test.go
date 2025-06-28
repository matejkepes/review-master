//go:build !integration
package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	main()
}

// update reviews (normal operation)
func TestProcess(t *testing.T) {
	Process(-1, false, false)
}

// report 2 months back from this month
func TestProcessReport(t *testing.T) {
	Process(2, false, false)
}

// report previous month
func TestProcessReport1(t *testing.T) {
	Process(1, false, false)
}

// report this month
func TestProcessReport2(t *testing.T) {
	Process(0, false, false)
}

// report previous month, send CSV report only
func TestProcessReport3(t *testing.T) {
	Process(1, false, true)
}

// report name or postal code not found
func TestProcessNameOrPostalCodeNotFoundReport1(t *testing.T) {
	Process(0, true, false)
}
