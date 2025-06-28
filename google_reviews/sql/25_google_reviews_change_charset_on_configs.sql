--
-- NOTE: This should only be run if updating an older database to change the configs charset to allow emojis
--
ALTER TABLE `google_reviews`.`google_reviews_configs` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
