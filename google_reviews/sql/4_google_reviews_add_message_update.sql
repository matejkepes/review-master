--
-- NOTE: This should only be run if updating an older database to use message from database
--
ALTER TABLE `google_reviews_configs` 
ADD COLUMN `use_database_message` TINYINT(1) NOT NULL DEFAULT 0 AFTER `multi_message_separator`,
ADD COLUMN `message` VARCHAR(2000) NOT NULL DEFAULT 'change me' AFTER `use_database_message`;
