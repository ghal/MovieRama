CREATE TABLE IF NOT EXISTS `users`
(
    `id`         int unsigned                                              NOT NULL AUTO_INCREMENT,
    `username`   varchar(255)                                              NOT NULL,
    `password`   char(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `first_name` varchar(255)                                              NOT NULL,
    `last_name`  varchar(255)                                              NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Store user details';
