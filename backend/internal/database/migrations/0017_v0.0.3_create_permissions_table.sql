-- +goose Up
-- +goose StatementBegin
CREATE TABLE `permissions` (`id` varchar(36),`name` text,`effect` text,`description` text,PRIMARY KEY (`id`));
CREATE UNIQUE INDEX `idx_permissions_name` ON `permissions`(`name`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `permissions`;
-- +goose StatementEnd
