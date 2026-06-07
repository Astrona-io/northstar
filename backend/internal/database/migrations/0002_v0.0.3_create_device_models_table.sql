-- +goose Up
-- +goose StatementBegin
CREATE TABLE `device_models` (`id` varchar(36),`manufacturer_id` varchar(36),`model_name` text,`general_info` text,`subtype` text,`revision` integer DEFAULT 1,`ports` text,`created_at` datetime,PRIMARY KEY (`id`),CONSTRAINT `fk_device_models_manufacturer` FOREIGN KEY (`manufacturer_id`) REFERENCES `manufacturers`(`id`));
CREATE INDEX `idx_device_models_model_name` ON `device_models`(`model_name`);
CREATE INDEX `idx_device_models_manufacturer_id` ON `device_models`(`manufacturer_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `device_models`;
-- +goose StatementEnd
