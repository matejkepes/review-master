--
-- NOTE: This should only be run if updating an older database to add the google_reviews_stats table
--

--
-- Table structure for table `google_reviews_stats`
--

DROP TABLE IF EXISTS `google_reviews`.`google_reviews_stats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `google_reviews`.`google_reviews_stats` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `stats_date` DATE NOT NULL DEFAULT '2000-01-01',
  `sent_count` int(10) unsigned NOT NULL DEFAULT '0',
  `requested_count` int(10) unsigned NOT NULL DEFAULT '0',
  `client_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `client_id_stats_date` (`client_id`,`stats_date`),
  INDEX `stats_date_ndx` (`stats_date` ASC)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
