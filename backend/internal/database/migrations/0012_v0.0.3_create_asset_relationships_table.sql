-- +goose Up
-- +goose StatementBegin
CREATE TABLE `asset_relationships` (`id` varchar(36),`source_id` text,`target_id` text,`type` text,PRIMARY KEY (`id`));
CREATE INDEX `idx_asset_relationships_target_id` ON `asset_relationships`(`target_id`);
CREATE INDEX `idx_asset_relationships_source_id` ON `asset_relationships`(`source_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `asset_relationships`;
-- +goose StatementEnd
