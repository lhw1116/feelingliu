-- create database
CREATE DATABASE
    IF NOT EXISTS `blog`;


-- create table `articles`
    CREATE TABLE `article` (
                               `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
                               `title` VARCHAR (32) NOT NULL,
                               `content` text NOT NULL,
                               `html` text NOT NULL,
                               `tag_id` INT UNSIGNED NOT NULL,
                               `created_time` VARCHAR (32) NOT NULL,
                               PRIMARY KEY (`id`),
                               UNIQUE KEY `title` (`title`) USING BTREE,
                               KEY `tags` (`tag_id`)
    ) ENGINE = INNODB AUTO_INCREMENT = 55 DEFAULT CHARSET = utf8mb4;

-- create table tag
CREATE TABLE `tag` (
                            `id` INT (10) UNSIGNED NOT NULL AUTO_INCREMENT,
                            `tag_name` VARCHAR (16) NOT NULL,
                            PRIMARY KEY (`id`)
) ENGINE = INNODB AUTO_INCREMENT = 6 DEFAULT CHARSET = utf8mb4;

--  create table user

CREATE TABLE `user` (
                        `id` INT NOT NULL AUTO_INCREMENT,
                        `username` VARCHAR (255) DEFAULT NULL,
                        `password` VARCHAR (255) DEFAULT NULL,
                        PRIMARY KEY (`id`)
) ENGINE = INNODB AUTO_INCREMENT = 2 DEFAULT CHARSET = utf8mb4;

--  insert user data

INSERT INTO `user` VALUES ('1', 'liuhanwen', '1111', '111', 'https://img2.woyaogexing.com/2020/04/14/c79f15db59c149368f8728a63f91346c!400x400.jpeg', 'Feelingliu', ' 111！\n');