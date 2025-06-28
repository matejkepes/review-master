--
-- NOTE: This should only be run if updating an older database to use dispatcher checks from database
--
ALTER TABLE `google_reviews_configs`
ADD COLUMN `dispatcher_checks_enabled` TINYINT(1) NOT NULL DEFAULT 0 AFTER `message`,
ADD COLUMN `booking_id_parameter` VARCHAR(255) NOT NULL DEFAULT 'b' AFTER `dispatcher_checks_enabled`,
ADD COLUMN `is_booking_for_now_diff_minutes` int(10) unsigned NOT NULL DEFAULT 10 AFTER `booking_id_parameter`,
ADD COLUMN `booking_now_pickup_to_contact_minutes` int(10) unsigned NOT NULL DEFAULT 10 AFTER `is_booking_for_now_diff_minutes`,
ADD COLUMN `pre_booking_pickup_to_contact_minutes` int(10) unsigned NOT NULL DEFAULT 3 AFTER `booking_now_pickup_to_contact_minutes`;
