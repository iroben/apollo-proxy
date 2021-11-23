CREATE TABLE `build_info`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `project`   varchar(32) DEFAULT NULL,
    `namespace` varchar(32) DEFAULT NULL,
    `job_id`    int         DEFAULT NULL,
    `result`    text,
    `created`   datetime    DEFAULT NULL,
    `state`     varchar(16) DEFAULT NULL,
    `branch`    varchar(32) DEFAULT NULL,
    `env`       varchar(32) DEFAULT NULL,
    `count`     tinyint     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `job_id` (`job_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;