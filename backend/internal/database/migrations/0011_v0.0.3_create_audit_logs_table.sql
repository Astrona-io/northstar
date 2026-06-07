-- +goose Up
-- +goose StatementBegin
CREATE TABLE `audit_logs` (`id` varchar(36),`asset_id` text,`action` text,`changes` text,`timestamp` datetime,PRIMARY KEY (`id`));
CREATE INDEX `idx_audit_logs_asset_id` ON `audit_logs`(`asset_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `audit_logs`;
-- +goose StatementEnd
