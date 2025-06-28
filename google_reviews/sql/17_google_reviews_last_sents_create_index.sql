--
-- NOTE: This should only be run if updating an older database to add index on google reviews last sents
-- on last_sent_date (which needs to be created) and client_id which is used a lot for counting daily sent and it is slow.
-- The last_sent_date is created because it is this that is used in the query. Previously a DATE(last_sent) was used
-- which is slow as it is not possible to create an index on this.
--
-- ALTER TABLE `google_reviews`.`google_reviews_last_sents` 
-- ADD INDEX `client_id_last_sent` (`last_sent` ASC, `client_id` ASC);
-- ALTER TABLE `google_reviews`.`google_reviews_last_sents` 
-- DROP INDEX `client_id_last_sent` ;

ALTER TABLE `google_reviews`.`google_reviews_last_sents` 
ADD COLUMN `last_sent_date` DATE NOT NULL DEFAULT '2000-01-01' AFTER `last_sent`;

-- NOTE: need to use a self join to update the last_sent_date based on last_sent
-- NOTE: if get Error Code: 1175. You are using safe update mode and you tried to update a table without a WHERE that uses a KEY column.  To disable safe mode, toggle the option in Preferences -> SQL Editor and reconnect.
--   SET SQL_SAFE_UPDATES = 0;
-- then set it back again after performing the update
-- SET SQL_SAFE_UPDATES = 1;
UPDATE `google_reviews`.`google_reviews_last_sents` AS l1, `google_reviews`.`google_reviews_last_sents` AS l2
SET l1.`last_sent_date` = DATE(l2.`last_sent`) WHERE l1.`id` = l2.`id`;

ALTER TABLE `google_reviews`.`google_reviews_last_sents` 
ADD INDEX `client_id_last_sent_date` (`last_sent_date` ASC, `client_id` ASC);
