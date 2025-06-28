--
-- NOTE: This should only be run if updating an older database to use dispatcher checks from database
--
ALTER TABLE `google_reviews`.`google_reviews_configs` 
ADD COLUMN `dispatcher_url` VARCHAR(255) NOT NULL DEFAULT '' AFTER `dispatcher_checks_enabled`;
