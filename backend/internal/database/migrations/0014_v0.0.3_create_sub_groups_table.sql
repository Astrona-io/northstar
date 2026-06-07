-- +goose Up
-- +goose StatementBegin
CREATE TABLE `sub_groups` (`id` varchar(36),`category_id` varchar(36),`name` text,`icon` text,`description` text,`item_order` integer DEFAULT 0,PRIMARY KEY (`id`),CONSTRAINT `fk_categories_sub_groups` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`) ON DELETE CASCADE);
CREATE INDEX `idx_sub_groups_category_id` ON `sub_groups`(`category_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `sub_groups`;
-- +goose StatementEnd
