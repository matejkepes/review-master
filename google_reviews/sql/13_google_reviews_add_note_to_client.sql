--
-- NOTE: This should only be run if updating an older database to dispatcher type
--
ALTER TABLE `google_reviews`.`clients` 
ADD COLUMN `note` VARCHAR(2000) NOT NULL AFTER `name`;
