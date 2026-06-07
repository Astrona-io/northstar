-- +goose Up
-- +goose StatementBegin
CREATE TABLE `datacenter_types` (`id` varchar(36),`name` text,PRIMARY KEY (`id`));
CREATE UNIQUE INDEX `idx_datacenter_types_name` ON `datacenter_types`(`name`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `datacenter_types`;
-- +goose StatementEnd
