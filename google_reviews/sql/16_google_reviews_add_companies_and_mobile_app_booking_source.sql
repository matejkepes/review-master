--
-- NOTE: This should only be run if updating an older database to add companies and booking source mobile app state info
--
ALTER TABLE `google_reviews`.`google_reviews_configs`
ADD COLUMN `companies` VARCHAR(255) NOT NULL DEFAULT '' AFTER `review_master_sms_gateway_pair_code`,
ADD COLUMN `booking_source_mobile_app_state` TINYINT(1) NOT NULL DEFAULT -1 AFTER `companies`
