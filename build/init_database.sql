CREATE DATABASE IF NOT EXISTS `booking`;

USE `booking`;

CREATE TABLE IF NOT EXISTS `screen` (
                          `id` int(11) NOT NULL AUTO_INCREMENT,
                          `number_seat_row` int(11) NOT NULL,
                          `number_seat_column` int(11) NOT NULL,
                          `created_date` datetime NOT NULL,
                          `user_id` varchar(20) NOT NULL,
                          PRIMARY KEY (`id`),
                          KEY `user_id_idx_screen_id` (`user_id`),
                          CONSTRAINT `screen_id_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS  `seat` (
                        `id` int(11) NOT NULL AUTO_INCREMENT,
                        `row_id` int(11) NOT NULL,
                        `column_id` int(11) NOT NULL,
                        `user_id` varchar(50) NOT NULL,
                        `screen_id` int(11) NOT NULL,
                        `booked_date` datetime NOT NULL,
                        PRIMARY KEY (`id`,`row_id`,`column_id`),
                        KEY `user_id_idx` (`user_id`),
                        KEY `id_idx` (`screen_id`),
                        CONSTRAINT `screen_id` FOREIGN KEY (`screen_id`) REFERENCES `screen` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
                        CONSTRAINT `user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=130 DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `user` (
                        `id` varchar(50) NOT NULL,
                        `password` varchar(200) NOT NULL,
                        PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `user` VALUES ('tester','$2a$14$BdgOuNVBU7sdGW9rIDIIv.MWXDdvTVKyTppb3qW03bmvz/6hhA1FO');
