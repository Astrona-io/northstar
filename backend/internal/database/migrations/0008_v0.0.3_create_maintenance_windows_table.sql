-- +goose Up
-- +goose StatementBegin
CREATE TABLE `maintenance_windows` (`id` varchar(36),`asset_id` varchar(36),`title` text,`description` text,`start_time` datetime,`end_time` datetime,`status` text DEFAULT "scheduled",PRIMARY KEY (`id`),CONSTRAINT `fk_assets_maintenance_windows` FOREIGN KEY (`asset_id`) REFERENCES `assets`(`id`) ON DELETE CASCADE);
CREATE INDEX `idx_maintenance_windows_title` ON `maintenance_windows`(`title`);
CREATE INDEX `idx_maintenance_windows_asset_id` ON `maintenance_windows`(`asset_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `maintenance_windows`;
-- +goose StatementEnd
