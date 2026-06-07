-- +goose Up
-- +goose StatementBegin
CREATE TABLE `device_model_revisions` (`id` varchar(36),`device_model_id` text,`revision` integer,`model_name` text,`general_info` text,`subtype` text,`ports` text,`created_at` datetime,PRIMARY KEY (`id`));
CREATE INDEX `idx_device_model_revisions_device_model_id` ON `device_model_revisions`(`device_model_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `device_model_revisions`;
-- +goose StatementEnd
