-- +goose Up
-- +goose StatementBegin
CREATE TABLE `webhooks` (`id` varchar(36),`url` text,`event` text,PRIMARY KEY (`id`));
CREATE INDEX `idx_webhooks_event` ON `webhooks`(`event`);
CREATE INDEX `idx_webhooks_url` ON `webhooks`(`url`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `webhooks`;
-- +goose StatementEnd
