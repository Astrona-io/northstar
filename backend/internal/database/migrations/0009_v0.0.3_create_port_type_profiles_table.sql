-- +goose Up
-- +goose StatementBegin
CREATE TABLE `port_type_profiles` (`id` varchar(36),`type` text,`name` text,`speeds` text,`created_at` datetime,PRIMARY KEY (`id`));
CREATE UNIQUE INDEX `idx_port_type_profiles_type` ON `port_type_profiles`(`type`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `port_type_profiles`;
-- +goose StatementEnd
