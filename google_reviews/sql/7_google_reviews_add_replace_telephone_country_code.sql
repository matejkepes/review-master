--
-- NOTE: This should only be run if updating an older database to use dispatcher checks from database
--
ALTER TABLE `google_reviews`.`google_reviews_configs` 
ADD COLUMN `replace_telephone_country_code` TINYINT(1) NOT NULL DEFAULT 0 AFTER `pre_booking_pickup_to_contact_minutes`,
ADD COLUMN `replace_telephone_country_code_with` VARCHAR(255) NOT NULL DEFAULT '0' AFTER `replace_telephone_country_code`;
