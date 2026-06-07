-- +goose Up
-- +goose StatementBegin
CREATE TABLE `datacenters` (`id` varchar(36),`name` text,`datacenter_type_id` varchar(36),`type` text,`country` text,`city` text,`properties` text,PRIMARY KEY (`id`),CONSTRAINT `fk_datacenters_datacenter_type` FOREIGN KEY (`datacenter_type_id`) REFERENCES `datacenter_types`(`id`));
CREATE INDEX `idx_datacenters_datacenter_type_id` ON `datacenters`(`datacenter_type_id`);
CREATE UNIQUE INDEX `idx_datacenters_name` ON `datacenters`(`name`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `datacenters`;
-- +goose StatementEnd
