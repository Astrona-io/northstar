-- +goose Up
-- +goose StatementBegin
CREATE TABLE `datacenter_walls` (`id` varchar(36),`datacenter_floor_id` varchar(36),`x1` real,`y1` real,`x2` real,`y2` real,`thickness` real DEFAULT 8,`type` text DEFAULT "wall",`flipped` numeric DEFAULT false,PRIMARY KEY (`id`),CONSTRAINT `fk_datacenter_floors_walls` FOREIGN KEY (`datacenter_floor_id`) REFERENCES `datacenter_floors`(`id`) ON DELETE CASCADE);
CREATE INDEX `idx_datacenter_walls_datacenter_floor_id` ON `datacenter_walls`(`datacenter_floor_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `datacenter_walls`;
-- +goose StatementEnd
