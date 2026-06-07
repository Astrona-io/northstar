-- +goose Up
-- +goose StatementBegin
CREATE TABLE `manufacturers` (`id` varchar(36),`name` text,`website` text,PRIMARY KEY (`id`));
CREATE UNIQUE INDEX `idx_manufacturers_name` ON `manufacturers`(`name`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `manufacturers`;
-- +goose StatementEnd
