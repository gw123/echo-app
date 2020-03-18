-- +goose Up
ALTER TABLE `echoapp`.`accounts` ADD INDEX `idx_accounts_meid` (`meid`);

-- +goose Down
ALTER TABLE `echoapp`.`accounts` DROP INDEX `idx_accounts_meid`;