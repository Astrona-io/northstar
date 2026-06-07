-- +goose Up
-- +goose StatementBegin
CREATE TABLE `custom_fields` (`id` varchar(36),`name` text,`field_type` text,`asset_type` text,`is_required` numeric DEFAULT false,`tab_group` text,`validation_regex` text,PRIMARY KEY (`id`));
CREATE INDEX `idx_custom_fields_asset_type` ON `custom_fields`(`asset_type`);
CREATE UNIQUE INDEX `idx_custom_fields_name` ON `custom_fields`(`name`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `custom_fields`;
-- +goose StatementEnd
