CREATE TABLE `project`
(
    `id`            int unsigned NOT NULL AUTO_INCREMENT,
    `hash`          varchar(32) DEFAULT NULL,
    `project`       varchar(32) DEFAULT NULL,
    `branch`        varchar(32) DEFAULT NULL,
    `apollo_app_id` varchar(32) DEFAULT NULL,
    `cluster_name`  varchar(32) DEFAULT NULL,
    `namespace`     varchar(32) DEFAULT NULL,
    `env`           varchar(32) DEFAULT NULL,
    `created`       datetime    DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `hash` (`hash`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;