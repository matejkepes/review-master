-- Add report_email_address field to clients table
-- This allows specifying a custom email address for sending reports to clients
ALTER TABLE
    `clients`
ADD
    COLUMN `report_email_address` VARCHAR(255) NULL DEFAULT NULL COMMENT 'Email address for sending reports to clients';