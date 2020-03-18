-- +goose Up
CREATE TABLE `echoapp`.`accounts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  `device_no` varchar(100) DEFAULT NULL COMMENT '设备号',
  `token` varchar(100) DEFAULT NULL COMMENT '身份唯一标识',
  `secret` varchar(100) DEFAULT NULL COMMENT '对称密钥',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_accounts_device_no` (`device_no`),
  UNIQUE KEY `uix_accounts_token` (`token`),
  KEY `idx_accounts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
DROP TABLE `echoapp`.`accounts`;