--
-- NOTE: This should only be run if updating an older database to add index on google reviews last sents
-- added to speed up stats query, using last_sent_date since quicker than using last_sent
-- also used for check nothing sent
-- The queries need changing and code modification to make these work faster

ALTER TABLE `google_reviews`.`google_reviews_last_sents` 
ADD INDEX `last_sent_date_ndx` (`last_sent_date` ASC);
