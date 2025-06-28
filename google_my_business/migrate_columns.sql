-- Add columns to google_reviews_configs table if they don't exist
ALTER TABLE
    google_reviews_configs
ADD
    COLUMN IF NOT EXISTS ai_responses_enabled BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE
    google_reviews_configs
ADD
    COLUMN IF NOT EXISTS contact_method VARCHAR(255) NULL;

-- Migrate data from clients to google_reviews_configs
UPDATE
    google_reviews_configs grc
    JOIN clients c ON c.id = grc.client_id
SET
    grc.ai_responses_enabled = COALESCE(c.ai_responses_enabled, FALSE),
    grc.contact_method = c.contact_method;

-- After verifying the migration, remove columns from clients table
-- ALTER TABLE clients DROP COLUMN ai_responses_enabled;
-- ALTER TABLE clients DROP COLUMN contact_method;