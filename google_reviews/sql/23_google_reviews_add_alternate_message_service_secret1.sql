--
-- NOTE: This should only be run if updating an older database to add companies and booking source mobile app state info
--
ALTER TABLE `google_reviews`.`google_reviews_configs`
ADD COLUMN `alternate_message_service_secret1` VARCHAR(255) NOT NULL DEFAULT '' AFTER `alternate_message_service`;
