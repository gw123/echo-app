-- +goose Up
ALTER TABLE `echoapp`.`upgrades`
ADD COLUMN `link` varchar(1024) DEFAULT '' COMMENT '文件公共下载链接';

-- +goose Down
ALTER TABLE `echoapp`.`upgrades`
DROP COLUMN `link`;