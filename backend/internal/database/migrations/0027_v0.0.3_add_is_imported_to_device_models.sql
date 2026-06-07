-- +goose Up
-- +goose StatementBegin
-- Add is_imported columns dynamically under version v0.0.3 release
ALTER TABLE `device_models` ADD COLUMN `is_imported` numeric DEFAULT false;
ALTER TABLE `device_model_revisions` ADD COLUMN `is_imported` numeric DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SQLite does not easily support dropping columns without recreating tables; we keep down blank for simplicity
-- +goose StatementEnd
