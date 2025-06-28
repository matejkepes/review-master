--
-- NOTE: This should only be run if updating an older database to use days for times
--
ALTER TABLE `google_reviews`.`google_reviews_config_times` 
ADD COLUMN `sunday` TINYINT(1) NOT NULL DEFAULT 1 AFTER `end`,
ADD COLUMN `monday` TINYINT(1) NOT NULL DEFAULT 1 AFTER `sunday`,
ADD COLUMN `tuesday` TINYINT(1) NOT NULL DEFAULT 1 AFTER `monday`,
ADD COLUMN `wednesday` TINYINT(1) NOT NULL DEFAULT 1 AFTER `tuesday`,
ADD COLUMN `thursday` TINYINT(1) NOT NULL DEFAULT 1 AFTER `wednesday`,
ADD COLUMN `friday` TINYINT(1) NOT NULL DEFAULT 1 AFTER `thursday`,
ADD COLUMN `saturday` TINYINT(1) NOT NULL DEFAULT 1 AFTER `friday`;
