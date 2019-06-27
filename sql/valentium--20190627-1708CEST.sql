#
# SQL Export
# Created by Querious (201054)
# Created: 27 June 2019 at 17:08:25 CEST
# Encoding: Unicode (UTF-8)
#


SET @PREVIOUS_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS;
SET FOREIGN_KEY_CHECKS = 0;


DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `accounts`;


CREATE TABLE `accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userId` varchar(128) DEFAULT NULL,
  `licenseId` int(11) DEFAULT '0',
  `expireDate` date DEFAULT NULL,
  `createdDate` datetime DEFAULT NULL,
  `updatedDate` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='\n';


CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ou` varchar(128) DEFAULT NULL,
  `oe` varbinary(255) DEFAULT NULL,
  `oh` varchar(128) DEFAULT NULL,
  `ca` varchar(128) DEFAULT NULL,
  `rn` varbinary(255) DEFAULT NULL,
  `bn` varbinary(255) DEFAULT NULL,
  `em` varbinary(255) DEFAULT NULL,
  `ci` int(11) DEFAULT '0',
  `si` int(11) DEFAULT '0',
  `ap` varbinary(255) DEFAULT NULL,
  `ie` varbinary(255) DEFAULT NULL,
  `li` int(11) DEFAULT '0',
  `xi` int(11) DEFAULT '0',
  `db` date DEFAULT NULL,
  `dl` datetime DEFAULT NULL,
  `dc` datetime DEFAULT NULL,
  `du` datetime DEFAULT NULL,
  `f1` int(11) DEFAULT '0',
  `f2` int(11) DEFAULT '0',
  `st` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=124 DEFAULT CHARSET=utf8 COMMENT='\n';




SET FOREIGN_KEY_CHECKS = @PREVIOUS_FOREIGN_KEY_CHECKS;


SET @PREVIOUS_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS;
SET FOREIGN_KEY_CHECKS = 0;


LOCK TABLES `accounts` WRITE;
ALTER TABLE `accounts` DISABLE KEYS;
ALTER TABLE `accounts` ENABLE KEYS;
UNLOCK TABLES;


LOCK TABLES `users` WRITE;
ALTER TABLE `users` DISABLE KEYS;
ALTER TABLE `users` ENABLE KEYS;
UNLOCK TABLES;




SET FOREIGN_KEY_CHECKS = @PREVIOUS_FOREIGN_KEY_CHECKS;


