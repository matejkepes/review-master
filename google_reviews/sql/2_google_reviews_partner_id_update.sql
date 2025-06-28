--
-- NOTE: This should only be run if updating an older database to include partner id associated with clients
--
ALTER TABLE `google_reviews`.`clients` 
ADD COLUMN `partner_id` BIGINT(20) NOT NULL AFTER `country`;
