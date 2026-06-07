-- +goose Up
-- +goose StatementBegin
CREATE TABLE `group_permissions` (`group_id` varchar(36),`permission_id` varchar(36),PRIMARY KEY (`group_id`,`permission_id`),CONSTRAINT `fk_group_permissions_group` FOREIGN KEY (`group_id`) REFERENCES `groups`(`id`),CONSTRAINT `fk_group_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions`(`id`));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `group_permissions`;
-- +goose StatementEnd
