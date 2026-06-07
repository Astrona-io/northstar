-- +goose Up
-- +goose StatementBegin
CREATE TABLE `license_agreements` (`id` varchar(36),`user_id` varchar(36),`version` text,`accepted_at` integer,PRIMARY KEY (`id`),CONSTRAINT `fk_license_agreements_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`));
CREATE INDEX `idx_license_agreements_user_id` ON `license_agreements`(`user_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `license_agreements`;
-- +goose StatementEnd
