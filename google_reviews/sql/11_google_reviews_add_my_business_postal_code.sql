--
-- NOTE: This should only be run if updating an older database to add postal_code to reply to google my business reviews
--
ALTER TABLE `google_reviews`.`google_reviews_configs` 
ADD COLUMN `google_my_business_postal_code` VARCHAR(255) NOT NULL AFTER `google_my_business_location_name`;
