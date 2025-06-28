--
-- NOTE: This should only be run if updating an older database to add companies and booking source mobile app state info
--
ALTER TABLE `google_reviews`.`google_reviews_configs`
ADD COLUMN `alternate_message_service_enabled` TINYINT(1) NOT NULL DEFAULT 0 AFTER `review_master_sms_gateway_pair_code`,
ADD COLUMN `alternate_message_service` VARCHAR(255) NOT NULL DEFAULT '' AFTER `alternate_message_service_enabled`
