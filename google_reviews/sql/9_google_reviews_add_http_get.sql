--
-- NOTE: This should only be run if updating an older database to use http get method
--
ALTER TABLE `google_reviews`.`google_reviews_configs` 
ADD COLUMN `http_get` TINYINT(1) NOT NULL DEFAULT 0 AFTER `send_url`;
