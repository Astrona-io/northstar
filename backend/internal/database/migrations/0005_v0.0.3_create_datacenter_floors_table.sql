-- +goose Up
-- +goose StatementBegin
CREATE TABLE `datacenter_floors` (`id` varchar(36),`datacenter_id` varchar(36),`name` text,`level` integer,`width` real,`depth` real,PRIMARY KEY (`id`),CONSTRAINT `fk_datacenters_floors` FOREIGN KEY (`datacenter_id`) REFERENCES `datacenters`(`id`) ON DELETE CASCADE);
CREATE INDEX `idx_datacenter_floors_datacenter_id` ON `datacenter_floors`(`datacenter_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `datacenter_floors`;
-- +goose StatementEnd
