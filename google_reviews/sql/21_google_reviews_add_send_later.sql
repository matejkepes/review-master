--
-- NOTE: This should only be run if updating an older database to add companies and booking source mobile app state info
--
ALTER TABLE `google_reviews`.`google_reviews_configs`
ADD COLUMN `send_delay_enabled` TINYINT(1) NOT NULL DEFAULT 0 AFTER `message`,
ADD COLUMN `send_delay` INT(10) NOT NULL DEFAULT 10 AFTER `send_delay_enabled`;

--
-- Table structure for table `google_reviews_send_laters`
--

DROP TABLE IF EXISTS `google_reviews`.`google_reviews_send_laters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `google_reviews`.`google_reviews_send_laters` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `telephone` VARCHAR(15) NOT NULL,
  `send_after` DATETIME NOT NULL,
  `send_url` VARCHAR(255) NOT NULL,
  `http_method` VARCHAR(255) NOT NULL,
  `app_key` VARCHAR(255) NOT NULL,
  `secret_key` VARCHAR(255) NOT NULL,
  `http_headers` VARBINARY(2000),
  `http_params` VARCHAR(2000),
  `http_body` VARBINARY(2000),
  `send_from_icabbi_app` TINYINT(1) NOT NULL DEFAULT 0,
  `review_master_sms_gateway_enabled` TINYINT(1) NOT NULL DEFAULT 0,
  `alternate_message_service_enabled` TINYINT(1) NOT NULL DEFAULT 0,
  `alternate_message_service` VARCHAR(255) NOT NULL DEFAULT '',
  `send_from_own_sms_gateway_enabled` TINYINT(1) NOT NULL DEFAULT 0,
  `send_success_response` VARCHAR(255) NOT NULL,
  `max_daily_send_count` int(10) NOT NULL,
  `client_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `client_id_telephone` (`client_id`,`telephone`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
