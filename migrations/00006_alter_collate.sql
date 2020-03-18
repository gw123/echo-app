-- +goose Up
ALTER TABLE `echoapp`.`accounts` CONVERT TO CHARACTER SET utf8 COLLATE utf8_unicode_ci;
ALTER TABLE `echoapp`.`register_records` CONVERT TO CHARACTER SET utf8 COLLATE utf8_unicode_ci;
ALTER TABLE `echoapp`.`upgrades` CONVERT TO CHARACTER SET utf8 COLLATE utf8_unicode_ci;

-- +goose Down
ALTER TABLE `echoapp`.`accounts` CONVERT TO CHARACTER SET utf8 COLLATE utf8_general_ci;
ALTER TABLE `echoapp`.`register_records` CONVERT TO CHARACTER SET utf8 COLLATE utf8_general_ci;
ALTER TABLE `echoapp`.`upgrades` CONVERT TO CHARACTER SET utf8 COLLATE utf8_general_ci;
