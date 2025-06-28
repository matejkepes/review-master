--
-- NOTE: This should only be run if updating an older database to reply to google my business reviews
--
ALTER TABLE `google_reviews`.`google_reviews_configs` 
ADD COLUMN `google_my_business_review_reply_enabled` TINYINT(1) NOT NULL DEFAULT 0 AFTER `replace_telephone_country_code_with`,
ADD COLUMN `google_my_business_location_name` VARCHAR(255) NOT NULL AFTER `google_my_business_review_reply_enabled`,
ADD COLUMN `google_my_business_reply_to_unspecfified_star_rating` TINYINT(1) NOT NULL DEFAULT 0 AFTER `google_my_business_location_name`,
ADD COLUMN `google_my_business_unspecfified_star_rating_reply` TEXT NOT NULL AFTER `google_my_business_reply_to_unspecfified_star_rating`,
ADD COLUMN `google_my_business_reply_to_one_star_rating` TINYINT(1) NOT NULL DEFAULT 0 AFTER `google_my_business_unspecfified_star_rating_reply`,
ADD COLUMN `google_my_business_one_star_rating_reply` TEXT NOT NULL AFTER `google_my_business_reply_to_one_star_rating`,
ADD COLUMN `google_my_business_reply_to_two_star_rating` TINYINT(1) NOT NULL DEFAULT 0 AFTER `google_my_business_one_star_rating_reply`,
ADD COLUMN `google_my_business_two_star_rating_reply` TEXT NOT NULL AFTER `google_my_business_reply_to_two_star_rating`,
ADD COLUMN `google_my_business_reply_to_three_star_rating` TINYINT(1) NOT NULL DEFAULT 0 AFTER `google_my_business_two_star_rating_reply`,
ADD COLUMN `google_my_business_three_star_rating_reply` TEXT NOT NULL AFTER `google_my_business_reply_to_three_star_rating`,
ADD COLUMN `google_my_business_reply_to_four_star_rating` TINYINT(1) NOT NULL DEFAULT 0 AFTER `google_my_business_three_star_rating_reply`,
ADD COLUMN `google_my_business_four_star_rating_reply` TEXT NOT NULL AFTER `google_my_business_reply_to_four_star_rating`,
ADD COLUMN `google_my_business_reply_to_five_star_rating` TINYINT(1) NOT NULL DEFAULT 0 AFTER `google_my_business_four_star_rating_reply`,
ADD COLUMN `google_my_business_five_star_rating_reply` TEXT NOT NULL AFTER `google_my_business_reply_to_five_star_rating`;
