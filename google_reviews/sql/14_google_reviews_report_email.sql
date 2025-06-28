--
-- NOTE: This should only be run if updating an older database to add report email
--
ALTER TABLE `google_reviews`.`google_reviews_configs` 
ADD COLUMN `google_my_business_report_enabled` TINYINT(1) NOT NULL DEFAULT 0 AFTER `google_my_business_five_star_rating_reply`,
ADD COLUMN `email_address` VARCHAR(255) NOT NULL DEFAULT '' AFTER `google_my_business_report_enabled`;
