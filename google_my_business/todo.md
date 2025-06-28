# Monthly Review Analysis Implementation Plan

## 🚨 WORKING GUIDELINES

**Before marking any task as [x] DONE, we must verify:**

### ✅ Testing Requirements

- **All tests must pass:** `go test -tags=integration ./...`
- **Integration tests included:** Ensure new functionality has appropriate test coverage

### 🖥️ Report/Template Testing

- **Preview server testing:** `cd cmd/template_preview && go run main.go`
- **Visual verification:** Check both HTML rendering and PDF generation
- **Cross-browser compatibility:** Verify template renders correctly
- **Print/PDF compatibility:** Test `@media print` functionality

### 📋 Verification Checklist

- [ ] Code compiles without errors
- [ ] All existing tests pass
- [ ] New tests written and passing
- [ ] Report preview server shows changes correctly
- [ ] PDF generation works with new features
- [ ] No regression in existing functionality

---

## 1. Core Data Fetching

- [x] Implement `FetchReviews` function to retrieve reviews for a specific location and time period
- [x] Add proper pagination support for fetching large numbers of reviews
- [x] Implement date filtering to focus on monthly reviews

## 2. Analysis Data Structure

- [x] Define input JSON schema for providing review data to the LLM
- [x] Define output JSON schema with detailed comment guidance
- [x] Define storage schema for saving analysis results

## 3. LLM Integration

- [x] Create a new `ReviewAnalyzer` in the `review` package that follows the existing architecture pattern
- [x] Design `ReviewBatch` struct to hold multiple reviews for collective analysis
- [x] Create analysis prompt template with embedded JSON schema and instructions
- [x] Add JSON validation and extraction to handle LLM responses

## 4. Analysis Storage & Processing Pipeline

### 4.1 Storage Infrastructure

- [x] Design database schema for storing analysis results
  - [x] Create tables for reports, locations, and metadata
  - [x] Define appropriate indexes for efficient retrieval
  - [x] Plan for historical data storage and comparison
- [x] Implement database access layer
  - [x] Create functions to save complete analysis results
  - [x] Add query functions for retrieving reports by client, location, and date

### 4.2 Processing Implementation

- [x] Adapt existing review fetching code for monthly analysis

  - [x] Create function to fetch all clients with review analysis enabled
  - [x] Reuse `FetchReviews` from google_my_business_api.go for retrieving reviews
  - [x] Add date range filtering specific to monthly reports

- [x] Develop the main processing pipeline

  - [x] Create the `AnalyzeClientReviews` orchestration function (skeleton implementation)
  - [x] Integrate review fetching within the orchestration flow
  - [x] Connect analyzer to process fetched reviews
  - [x] Save analysis results to the database
  - [x] Add progress tracking and logging

- [x] Implement report generation and delivery
  - [x] Create HTML templates for reports
  - [x] Convert HTML reports to PDF using third-party service
  - [x] Add email delivery mechanism for reports
    - [x] Create email service
    - [x] Use email service to send reports to clients

## 5. Automation & Scheduling

- [x] Create a monthly scheduling mechanism
  - [x] Implement cron job that runs 2 AM on the 1st of each month (UK time)
  - [x] Create monthly_analysis.sh shell script
  - [x] Update main.go to support command-line arguments for monthly analysis
- [x] Implement error handling, retry logic and fallbacks
  - [x] Create file-based retry tracking (JSON) with 3 retry limit
  - [x] Support for -retry-only mode to process just failed clients
  - [x] Log errors but continue processing other clients
- [x] Add monitoring and logging for processing status
  - [x] Enhance ProcessingSummary with more detailed error information
  - [x] Create daily log files for each run
  - [x] Clean up old logs (older than 3 months)
- [x] Develop notification system for completed reports
  - [x] Send summary email to matejkepes@gmail.com after processing
  - [x] Include statistics on clients, reviews, PDFs, and failures

## 6. Presentation & Delivery

- [x] Design email templates for report delivery
- [x] Create dashboard views for the admin portal (if needed)
- [x] Implement filtering and sorting capabilities for reports
- [x] Add historical comparison features

## 7. Deployment & Documentation

- [x] Document the system architecture and implementation details
- [x] Update memory.md with details of new functionality
- [x] Create user documentation for interpreting reports
- [ ] Perform final testing with real-world data

## 8. Report Enhancement - Visual & Content Improvements

### 8.1 Sentiment Analysis Visualization Enhancement

- [x] Replace sentiment bar charts with pie charts
- [x] Use same colors as existing bars (green, blue, red)
- [x] Make pie chart larger (220px) and center horizontally
- [x] Position legend below the pie chart in horizontal layout
- [x] Ensure percentages always add up to 100% after rounding
- [x] Test with preview server and verify PDF compatibility

### 8.2 Negative Review Breakdown Visualization

- [x] Add pie chart for negative review categories in Overview section
  - [x] Create SVG-based pie chart for `NegativeReviewBreakdown.Categories` data
  - [x] Display category names and percentages
  - [x] Position section before "Suggested Improvements" in location card flow
  - [x] Use distinct color palette (8 colors) for different negative review categories
  - [x] Add section title "Negative Review Breakdown"
  - [x] Fixed visual pie chart discrepancy using original percentages for angles and rounded for display
  - [x] Fixed index out of range error by ensuring category count doesn't exceed available categories
  - [x] Implemented coherent test data generation that correlates with negative themes and improvement areas
  - [x] **COMPLETED**: All tests pass (`go test -tags=integration ./...`) and preview server verified ✅

### 8.3 Training Recommendations Section

- [x] Add "Training Recommendations" section after "Areas for Improvement"
  - [x] Create subsections for "For Operators" and "For Drivers"
  - [x] Use bullet list format for each subsection
  - [x] Add conditional display logic - hide section if no training data exists
  - [x] Style subsection headers with appropriate typography
  - [x] Generate location-specific training content that correlates with identified issues
  - [x] Add CSS styling for training-recommendations-container, subsections, and lists
  - [x] Include in print styles for proper page break handling
  - [x] **COMPLETED**: Training recommendations section successfully displays with location-specific content that correlates with negative themes and improvement areas. Proper conditional rendering and styling implemented. All tests pass and preview server works correctly.

### 8.4 Enhanced Suggestions Section

- [x] **SKIPPED** - Current "Suggested Improvements" section is sufficient as-is. No merge needed since "Actionable Fixes" section doesn't exist.

### 8.5 SVG Implementation and Styling

- [x] Implement reusable SVG pie chart generation functions
  - [x] Create `createPieChart()` JavaScript function for sentiment analysis
  - [x] Create `createNegativeBreakdownPieChart()` JavaScript function for negative categories
  - [x] Ensure SVG charts are PDF-compatible
  - [x] Add responsive sizing for different screen/print sizes
  - [x] Include proper labels and legends for accessibility
  - [x] **COMPLETED**: All SVG pie chart functions implemented and working correctly. Charts are PDF-compatible with proper labeling and responsive design.

### 8.6 CSS and Layout Updates

- [ ] Update location card layout to accommodate new sections
  - [ ] Ensure proper spacing between sections
  - [ ] Maintain page break functionality for PDF generation
  - [ ] Test responsive design on different screen sizes
  - [ ] Verify PDF rendering maintains quality of SVG charts

### 8.7 Testing and Validation

- [ ] Test HTML rendering in browser
- [ ] Test PDF generation with new pie charts
- [ ] Verify all existing data still displays correctly
- [ ] Test with various data scenarios (edge cases)
- [ ] Ensure training recommendations data populates correctly

## 9. First Page Enhancement - Key Metrics Dashboard

### 9.1 Executive Summary Metrics

- [x] Add key metrics dashboard section above location cards
  - [x] Display total reviews analyzed (sum across all locations)
  - [x] Show overall average rating (weighted by review count per location)
  - [x] Calculate and display overall positive sentiment percentage
  - [x] Design visually prominent metric cards with large numbers
  - [x] Add proper spacing and visual hierarchy
  - [x] Ensure responsive design and PDF compatibility
  - [x] Fix JavaScript calculations to extract data from DOM elements
  - [x] Add page break after executive summary for proper PDF pagination

### 9.2 Testing and Validation

- [x] Test HTML rendering in browser
- [x] Test PDF generation with new metrics dashboard
- [x] Verify calculations work correctly with real data
- [x] Ensure proper page breaks in print mode (first page: header + summary only)

### 9.3 Review Analysis Accuracy Fix ✅

- [x] **CRITICAL FIX**: Corrected review count discrepancy where reports showed only reviews with text content instead of total reviews
- [x] Updated analyzer to preserve total review count (including reviews without text) after AI analysis
- [x] **CONSISTENCY FIX**: Fixed average rating calculation to always use ALL reviews instead of conditionally using AI's calculation (which only included text reviews)
- [x] Added comprehensive test coverage for mixed review scenarios (TestAnalyzeWithMixedReviews)
- [x] Verified average ratings include all reviews regardless of text content
- [x] All tests pass and functionality verified with preview server
