--
-- NOTE: This should only be run if updating an older database to add companies and booking source mobile app state info
--
ALTER TABLE `google_reviews`.`google_reviews_configs`
ADD COLUMN `review_master_sms_gateway_use_master_queue` TINYINT(1) NOT NULL DEFAULT 0 AFTER `review_master_sms_gateway_enabled`;

--
-- Table structure for table `review_master_sms_gateway_master_queues`
--

DROP TABLE IF EXISTS `google_reviews`.`review_master_sms_gateway_master_queues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `google_reviews`.`review_master_sms_gateway_master_queues` (
  `id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

INSERT INTO `google_reviews`.`review_master_sms_gateway_master_queues` (`id`) VALUES (999999);
