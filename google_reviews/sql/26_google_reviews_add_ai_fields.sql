-- Migration: Add AI response configuration to clients table
ALTER TABLE
    google_reviews_configs
ADD
    COLUMN ai_responses_enabled BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE
    google_reviews_configs
ADD
    COLUMN contact_method VARCHAR(255) NULL;