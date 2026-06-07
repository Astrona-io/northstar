-- +goose Up
-- +goose StatementBegin
CREATE TABLE `network_interfaces` (`id` varchar(36),`asset_id` varchar(36),`nic_name` text,`name` text,`type` text,`mac_address` text,`ipv4_address` text,`ipv6_address` text,`speed` text,`mtu` integer DEFAULT 1500,`vlan` text,`status` text DEFAULT "up",PRIMARY KEY (`id`),CONSTRAINT `fk_assets_network_interfaces` FOREIGN KEY (`asset_id`) REFERENCES `assets`(`id`) ON DELETE CASCADE);
CREATE INDEX `idx_network_interfaces_name` ON `network_interfaces`(`name`);
CREATE INDEX `idx_network_interfaces_nic_name` ON `network_interfaces`(`nic_name`);
CREATE INDEX `idx_network_interfaces_asset_id` ON `network_interfaces`(`asset_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `network_interfaces`;
-- +goose StatementEnd
