CREATE TABLE IF NOT EXISTS `movies`
(
    `id`          int unsigned                                                  NOT NULL AUTO_INCREMENT,
    `title`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `user_id`     int unsigned                                                  NOT NULL,
    `description` text                                                          NOT NULL,
    `created_at`  timestamp                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `user_id` (`user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Store movie details';
