--
-- NOTE: This should only be run if updating an older database to send from iCabbi app
--
ALTER TABLE `google_reviews`.`google_reviews_configs` 
ADD COLUMN `send_from_icabbi_app` TINYINT(1) NOT NULL DEFAULT 0 AFTER `telephone_parameter`,
ADD COLUMN `app_key` VARCHAR(255) NOT NULL DEFAULT '' AFTER `send_from_icabbi_app`,
ADD COLUMN `secret_key` VARCHAR(255) NOT NULL DEFAULT '' AFTER `app_key`;
