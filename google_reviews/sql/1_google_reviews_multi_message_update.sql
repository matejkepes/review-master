--
-- NOTE: This should only be run if updating an older database to allow multi message
--
ALTER TABLE `google_reviews_configs` 
ADD COLUMN `multi_message_enabled` TINYINT(1) NOT NULL DEFAULT 0 AFTER `time_zone`,
ADD COLUMN `message_parameter` VARCHAR(255) NOT NULL DEFAULT 'm' AFTER `multi_message_enabled`,
ADD COLUMN `multi_message_separator` VARCHAR(255) NOT NULL DEFAULT 'SSSSS' AFTER `message_parameter`;
