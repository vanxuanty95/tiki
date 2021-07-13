CREATE DATABASE IF NOT EXISTS `booking`;

USE `booking`;

CREATE TABLE IF NOT EXISTS `screen` (
 `id` INT NOT NULL AUTO_INCREMENT,
 `number_seat_row` INT NOT NULL,
 `number_seat_column` VARCHAR(45) NOT NULL,
    PRIMARY KEY (`id`));
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8
    ROW_FORMAT = DYNAMIC;

CREATE TABLE IF NOT EXISTS `user` (
    `id` VARCHAR(50) NOT NULL,
    `password` VARCHAR(200) NOT NULL,
    PRIMARY KEY (`id`));
    ) ENGINE = InnoDB
    DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;


CREATE TABLE IF NOT EXISTS `seat` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `row` INT NOT NULL,
  `column` INT NOT NULL,
  `user_id` VARCHAR(50) NOT NULL,
  `screen_id` INT NOT NULL,
  `booked_date` DATETIME NOT NULL,
  PRIMARY KEY (`id`, `row`, `column`),
  INDEX `user_id_idx` (`user_id` ASC),
  INDEX `id_idx` (`screen_id` ASC),
  CONSTRAINT `user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `booking`.`user` (`id`)
  CONSTRAINT `screen_id`
    FOREIGN KEY (`screen_id`)
    REFERENCES `booking`.`screen` (`id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;