-- +goose Up
CREATE TABLE `coupons`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `com_id`      int                       default 0,
    `type`        varchar(12)      not null,
    `cover`       varchar(12)               default '',
    `name`        varchar(32)      not null,
    `desc`        varchar(256)     not null default '',
    `href`        varchar(128)              default '',
    `min_consume` int                       default 0,
    `amount`      decimal(6, 2)             default 0,
    `range_type`  char(6)                   default 'all' COMMENT 'all,range',
    `range`       varchar(1024)             default '[]',
    `total`       int                       default 0,
    `used_total`  int                       default 0,
    `expire_at`   timestamp        NULL     default NULL,
    `duration`    int                       default 0,
    `start_at`    timestamp        null     default null,
    `created_at`  timestamp        NULL     DEFAULT NULL COMMENT '创建时间',
    `updated_at`  timestamp        NULL     DEFAULT NULL COMMENT '更新时间',
    `deleted_at`  timestamp        NULL     DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

-- +goose Down
DROP TABLE `coupons`;