CREATE TABLE IF NOT EXISTS `movies_users_actions`
(
    `id`       int unsigned                                                          NOT NULL AUTO_INCREMENT,
    `user_id`  int unsigned                                                          NOT NULL,
    `movie_id` int unsigned                                                          NOT NULL,
    `action`   enum ('like','hate') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_id` (`user_id`, `movie_id`),
    UNIQUE KEY `action` (`action`, `movie_id`, `user_id`),
    KEY `user_id_2` (`user_id`),
    KEY `movie_id` (`movie_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Store movie user action details';
