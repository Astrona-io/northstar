-- +goose Up
-- +goose StatementBegin
CREATE TABLE `groups` (`id` varchar(36),`name` text,`description` text,PRIMARY KEY (`id`));
CREATE UNIQUE INDEX `idx_groups_name` ON `groups`(`name`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `groups`;
-- +goose StatementEnd
