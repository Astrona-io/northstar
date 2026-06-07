-- +goose Up
-- +goose StatementBegin
CREATE TABLE `users` (`id` varchar(36),`username` text,`password` text,`role` text,`group_id` varchar(36),PRIMARY KEY (`id`),CONSTRAINT `fk_users_group` FOREIGN KEY (`group_id`) REFERENCES `groups`(`id`));
CREATE INDEX `idx_users_group_id` ON `users`(`group_id`);
CREATE UNIQUE INDEX `idx_users_username` ON `users`(`username`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `users`;
-- +goose StatementEnd
