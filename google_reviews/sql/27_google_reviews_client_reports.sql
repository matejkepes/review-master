-- Migration: Create client_reports table for monthly review analysis
-- Description: Establishes the storage structure for client-level monthly review reports
CREATE TABLE client_reports (
    report_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    client_id INT NOT NULL,
    report_period_start DATE NOT NULL,
    report_period_end DATE NOT NULL,
    generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    locations JSON NOT NULL,
    -- Add unique constraint to prevent duplicate reports for the same client/period
    CONSTRAINT unique_client_period UNIQUE (
        client_id,
        report_period_start,
        report_period_end
    )
);

-- Add index for client_id lookups
CREATE INDEX idx_client_reports_client_id ON client_reports (client_id);

-- Add index for period-based queries
CREATE INDEX idx_client_reports_period ON client_reports (report_period_start, report_period_end);

-- Add comments to describe the table purpose
ALTER TABLE
    client_reports COMMENT = 'Stores monthly review analysis reports by client with location-level details stored as JSON';

-- monthly_review_analysis_enabled
ALTER TABLE
    google_reviews_configs
ADD
    COLUMN monthly_review_analysis_enabled BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Whether monthly review analysis is enabled for this client';

-- Add an index to improve query performance when filtering by this column
CREATE INDEX idx_monthly_review_analysis_enabled ON google_reviews_configs (monthly_review_analysis_enabled);