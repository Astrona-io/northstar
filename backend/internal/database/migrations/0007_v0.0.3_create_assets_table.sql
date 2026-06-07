-- +goose Up
-- +goose StatementBegin
CREATE TABLE `assets` (`id` varchar(36),`name` text,`type` text,`status` text DEFAULT "active",`ip_address` text,`description` text,`properties` text,`rack_id` varchar(36),`rack_position_u` integer,`device_model_id` varchar(36),`device_model_revision` integer,`host_asset_id` varchar(36),`container_image` text,`container_port_mapping` text,PRIMARY KEY (`id`),CONSTRAINT `fk_racks_assets` FOREIGN KEY (`rack_id`) REFERENCES `racks`(`id`),CONSTRAINT `fk_assets_device_model` FOREIGN KEY (`device_model_id`) REFERENCES `device_models`(`id`),CONSTRAINT `fk_assets_host_asset` FOREIGN KEY (`host_asset_id`) REFERENCES `assets`(`id`));
CREATE INDEX `idx_assets_host_asset_id` ON `assets`(`host_asset_id`);
CREATE INDEX `idx_assets_device_model_id` ON `assets`(`device_model_id`);
CREATE INDEX `idx_assets_rack_id` ON `assets`(`rack_id`);
CREATE INDEX `idx_assets_ip_address` ON `assets`(`ip_address`);
CREATE INDEX `idx_assets_type` ON `assets`(`type`);
CREATE INDEX `idx_assets_name` ON `assets`(`name`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `assets`;
-- +goose StatementEnd
