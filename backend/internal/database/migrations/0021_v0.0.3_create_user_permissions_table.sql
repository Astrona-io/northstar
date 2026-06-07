-- +goose Up
-- +goose StatementBegin
CREATE TABLE `user_permissions` (`user_id` varchar(36),`permission_id` varchar(36),PRIMARY KEY (`user_id`,`permission_id`),CONSTRAINT `fk_user_permissions_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),CONSTRAINT `fk_user_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions`(`id`));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `user_permissions`;
-- +goose StatementEnd
