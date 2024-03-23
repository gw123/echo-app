-- +goose Up
CREATE TABLE `user_coupons`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `com_id`     int                  default 0,
    `coupon_id`  int             not null,
    `user_id`    int             not null,
    `expire_at`  timestamp       NULL default NULL,
    `start_at`   timestamp       null default null,
    `created_at` timestamp       NULL DEFAULT NULL COMMENT '创建时间',
    `updated_at` timestamp       NULL DEFAULT NULL COMMENT '更新时间',
    `deleted_at` timestamp       NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

-- +goose Down
DROP TABLE `user_coupons`