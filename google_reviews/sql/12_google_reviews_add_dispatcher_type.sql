--
-- NOTE: This should only be run if updating an older database to dispatcher type
--
ALTER TABLE `google_reviews`.`google_reviews_configs` 
ADD COLUMN `dispatcher_type` VARCHAR(255) NOT NULL AFTER `dispatcher_url`;
