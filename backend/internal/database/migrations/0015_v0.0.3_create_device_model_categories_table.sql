-- +goose Up
-- +goose StatementBegin
CREATE TABLE `device_model_categories` (`device_model_id` varchar(36),`category_id` varchar(36),PRIMARY KEY (`device_model_id`,`category_id`),CONSTRAINT `fk_device_model_categories_device_model` FOREIGN KEY (`device_model_id`) REFERENCES `device_models`(`id`),CONSTRAINT `fk_device_model_categories_category` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `device_model_categories`;
-- +goose StatementEnd
