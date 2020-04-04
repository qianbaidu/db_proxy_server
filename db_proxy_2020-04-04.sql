# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.6.36)
# Database: db_proxy
# Generation Time: 2020-04-04 12:41:52 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table columns_permission
# ------------------------------------------------------------

DROP TABLE IF EXISTS `columns_permission`;

CREATE TABLE `columns_permission` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `table_perm_id` bigint(20) NOT NULL,
  `column_name` varchar(100) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `db_perm_id` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `db_id` (`db_perm_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `columns_permission` WRITE;
/*!40000 ALTER TABLE `columns_permission` DISABLE KEYS */;

INSERT INTO `columns_permission` (`id`, `table_perm_id`, `column_name`, `status`, `db_perm_id`)
VALUES
	(9,1,'mysql_password',0,1),
	(10,3,'mysql_password',0,3),
	(11,1,'mongodb_password',0,1),
	(12,3,'mongodb_password',0,3);

/*!40000 ALTER TABLE `columns_permission` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table db_list
# ------------------------------------------------------------

DROP TABLE IF EXISTS `db_list`;

CREATE TABLE `db_list` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL,
  `ip` varchar(64) NOT NULL,
  `port` int(5) NOT NULL,
  `db_name` varchar(32) NOT NULL,
  `username` varchar(32) NOT NULL,
  `password` varchar(128) NOT NULL,
  `encode` int(11) NOT NULL,
  `status` tinyint(1) NOT NULL COMMENT '是否有效，默认0有效',
  `create_time` datetime(6) NOT NULL,
  `db_type` int(11) NOT NULL COMMENT '1 mysql ',
  PRIMARY KEY (`id`),
  UNIQUE KEY `db` (`ip`,`port`,`db_name`),
  UNIQUE KEY `db_name` (`db_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `db_list` WRITE;
/*!40000 ALTER TABLE `db_list` DISABLE KEYS */;

INSERT INTO `db_list` (`id`, `name`, `ip`, `port`, `db_name`, `username`, `password`, `encode`, `status`, `create_time`, `db_type`)
VALUES
	(1,'db_proxy','127.0.0.1',3306,'db_proxy','root','',1,0,'2017-10-10 08:53:40.342052',1),
	(2,'test','127.0.0.1',27017,'test','admin','admin',1,0,'2017-10-10 08:53:40.342052',2);

/*!40000 ALTER TABLE `db_list` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table db_permission
# ------------------------------------------------------------

DROP TABLE IF EXISTS `db_permission`;

CREATE TABLE `db_permission` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `have_all` tinyint(1) NOT NULL DEFAULT '0',
  `have_secret_columns` tinyint(1) NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `db_id` int(11) NOT NULL DEFAULT '0',
  `user_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `db` (`db_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `db_permission` WRITE;
/*!40000 ALTER TABLE `db_permission` DISABLE KEYS */;

INSERT INTO `db_permission` (`id`, `have_all`, `have_secret_columns`, `status`, `db_id`, `user_id`)
VALUES
	(1,1,1,0,1,1),
	(2,1,0,0,2,1),
	(3,1,1,0,1,2);

/*!40000 ALTER TABLE `db_permission` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table event_log
# ------------------------------------------------------------

DROP TABLE IF EXISTS `event_log`;

CREATE TABLE `event_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `sql` longtext NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `create_time` datetime(6) NOT NULL,
  `db_id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `update_datatime` datetime(6) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `event_log` WRITE;
/*!40000 ALTER TABLE `event_log` DISABLE KEYS */;

INSERT INTO `event_log` (`id`, `sql`, `status`, `create_time`, `db_id`, `user_id`, `update_datatime`)
VALUES
	(1,'show databases',2,'2020-04-04 18:13:53.000000',1,2,'2020-04-04 18:13:53.000000'),
	(2,'show tables',2,'2020-04-04 18:13:53.000000',1,2,'2020-04-04 18:13:53.000000'),
	(3,'show databases',2,'2020-04-04 18:14:22.000000',1,2,'2020-04-04 18:14:22.000000'),
	(4,'show tables',2,'2020-04-04 18:14:32.000000',1,2,'2020-04-04 18:14:32.000000'),
	(5,'show databases',2,'2020-04-04 18:27:44.000000',1,2,'2020-04-04 18:27:44.000000'),
	(6,'show tables',2,'2020-04-04 18:27:44.000000',1,2,'2020-04-04 18:27:44.000000'),
	(7,'show databases',2,'2020-04-04 18:27:53.000000',1,2,'2020-04-04 18:27:53.000000'),
	(8,'show tables',2,'2020-04-04 18:27:59.000000',1,2,'2020-04-04 18:27:59.000000'),
	(9,'show tables',2,'2020-04-04 18:28:12.000000',1,2,'2020-04-04 18:28:12.000000'),
	(10,'show databases',2,'2020-04-04 18:29:21.000000',1,2,'2020-04-04 18:29:21.000000'),
	(11,'show tables',2,'2020-04-04 18:29:21.000000',1,2,'2020-04-04 18:29:21.000000'),
	(12,'show tables',2,'2020-04-04 18:29:23.000000',1,2,'2020-04-04 18:29:23.000000'),
	(13,'show databases',2,'2020-04-04 18:30:04.000000',1,2,'2020-04-04 18:30:04.000000'),
	(14,'show tables',2,'2020-04-04 18:30:04.000000',1,2,'2020-04-04 18:30:04.000000'),
	(15,'show tables',2,'2020-04-04 18:30:08.000000',1,2,'2020-04-04 18:30:08.000000'),
	(16,'show databases',2,'2020-04-04 18:30:44.000000',1,2,'2020-04-04 18:30:44.000000'),
	(17,'show tables',2,'2020-04-04 18:30:44.000000',1,2,'2020-04-04 18:30:44.000000'),
	(18,'show databases',2,'2020-04-04 18:36:23.000000',1,2,'2020-04-04 18:36:23.000000'),
	(19,'show databases',2,'2020-04-04 18:38:43.000000',1,2,'2020-04-04 18:38:43.000000'),
	(20,'show databases',2,'2020-04-04 18:42:50.000000',1,2,'2020-04-04 18:42:50.000000'),
	(21,'show databases',2,'2020-04-04 19:31:21.000000',1,2,'2020-04-04 19:31:21.000000'),
	(22,'show databases',2,'2020-04-04 19:31:25.000000',1,2,'2020-04-04 19:31:25.000000'),
	(23,'show databases',2,'2020-04-04 19:32:26.000000',1,2,'2020-04-04 19:32:26.000000'),
	(24,'show databases',2,'2020-04-04 19:40:46.000000',1,2,'2020-04-04 19:40:46.000000'),
	(25,'show databases',2,'2020-04-04 19:40:46.000000',1,2,'2020-04-04 19:40:46.000000'),
	(26,'show databases',2,'2020-04-04 19:40:47.000000',1,2,'2020-04-04 19:40:47.000000');

/*!40000 ALTER TABLE `event_log` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table table_permission
# ------------------------------------------------------------

DROP TABLE IF EXISTS `table_permission`;

CREATE TABLE `table_permission` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `table_name` varchar(100) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '默认0有效；1无效',
  `db_perm_id` int(11) NOT NULL DEFAULT '0',
  `have_secret` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `db_perm_id` (`db_perm_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `table_permission` WRITE;
/*!40000 ALTER TABLE `table_permission` DISABLE KEYS */;

INSERT INTO `table_permission` (`id`, `table_name`, `status`, `db_perm_id`, `have_secret`)
VALUES
	(1,'user',0,1,1),
	(3,'user',0,3,1);

/*!40000 ALTER TABLE `table_permission` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user` varchar(32) COLLATE utf8_bin NOT NULL DEFAULT '',
  `mysql_password` varchar(41) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `mysql_read_permission` tinyint(1) NOT NULL DEFAULT '0',
  `mysql_read_write_permission` tinyint(1) NOT NULL DEFAULT '0',
  `mongodb_password` varchar(41) COLLATE utf8_bin NOT NULL DEFAULT '',
  `mongodb_read_permission` tinyint(1) NOT NULL DEFAULT '0',
  `mongodb_read_write_permission` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user` (`user`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`id`, `user`, `mysql_password`, `mysql_read_permission`, `mysql_read_write_permission`, `mongodb_password`, `mongodb_read_permission`, `mongodb_read_write_permission`)
VALUES
	(1,X'616C657831','*8258F2618980E77E5220ECD738182656223809C1',1,0,X'3532353666323832636635323664613061636561613434386333343264313435',1,0),
	(2,X'74657374','*94BDCEBE19083CE2A1F959FD02F964C7AF4CFC29',1,0,X'6136646535323161626566633266656434663538373638353561333438346635',1,0);

/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
