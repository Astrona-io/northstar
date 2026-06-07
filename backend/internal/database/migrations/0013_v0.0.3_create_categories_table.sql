-- +goose Up
-- +goose StatementBegin
CREATE TABLE `categories` (`id` varchar(36),`name` text,`icon` text,`description` text,`item_order` integer DEFAULT 0,`parent_dependency` text,PRIMARY KEY (`id`));
CREATE UNIQUE INDEX `idx_categories_name` ON `categories`(`name`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `categories`;
-- +goose StatementEnd
