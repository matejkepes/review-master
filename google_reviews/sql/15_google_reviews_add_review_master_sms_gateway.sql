--
-- NOTE: This should only be run if updating an older database to add Review Master SMS Gateway info
--
ALTER TABLE `google_reviews`.`google_reviews_configs`
ADD COLUMN `review_master_sms_gateway_enabled` TINYINT(1) NOT NULL DEFAULT 0 AFTER `replace_telephone_country_code_with`,
ADD COLUMN `review_master_sms_gateway_pair_code` VARCHAR(255) NOT NULL DEFAULT '' AFTER `review_master_sms_gateway_enabled`;
-- ADD UNIQUE INDEX `review_master_sms_gateway_pair_code_UNIQUE` (`review_master_sms_gateway_pair_code` ASC);
