-- +goose Up
-- +goose StatementBegin
CREATE TABLE `racks` (`id` varchar(36),`datacenter_id` varchar(36),`datacenter_floor_id` varchar(36),`name` text,`height_u` integer DEFAULT 42,`placement_zone` text,`x` real,`y` real,PRIMARY KEY (`id`),CONSTRAINT `fk_datacenter_floors_racks` FOREIGN KEY (`datacenter_floor_id`) REFERENCES `datacenter_floors`(`id`) ON DELETE CASCADE,CONSTRAINT `fk_datacenters_racks` FOREIGN KEY (`datacenter_id`) REFERENCES `datacenters`(`id`) ON DELETE CASCADE);
CREATE INDEX `idx_racks_datacenter_floor_id` ON `racks`(`datacenter_floor_id`);
CREATE INDEX `idx_racks_datacenter_id` ON `racks`(`datacenter_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `racks`;
-- +goose StatementEnd
