-- +goose Up
-- +goose StatementBegin
-- Seed default custom fields
INSERT OR IGNORE INTO `custom_fields` (`id`, `name`, `field_type`, `asset_type`, `is_required`, `tab_group`)
VALUES ('warranty-expr-uuid', 'warranty_expiry', 'text', 'Server', 0, 'General Custom');

-- Seed default webhook
INSERT OR IGNORE INTO `webhooks` (`id`, `url`, `event`)
VALUES ('default-webhook-uuid', 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX', 'asset:create');

-- Seed default datacenter types
INSERT OR IGNORE INTO `datacenter_types` (`id`, `name`)
VALUES 
('dc-type-1', 'on-prem'),
('dc-type-2', 'homelab'),
('dc-type-3', 'air-gap'),
('dc-type-4', 'cloud:aws'),
('dc-type-5', 'cloud:gcp'),
('dc-type-6', 'cloud:azure');

-- Seed global default fallback category
INSERT OR IGNORE INTO `categories` (`id`, `name`, `icon`, `description`)
VALUES ('global-cat-uuid', 'global', 'i-heroicons-rectangle-group', 'Global default fallback hardware category');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM `custom_fields` WHERE `name` = 'warranty_expiry';
DELETE FROM `webhooks` WHERE `url` = 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX';
DELETE FROM `datacenter_types` WHERE `name` IN ('on-prem', 'homelab', 'air-gap', 'cloud:aws', 'cloud:gcp', 'cloud:azure');
DELETE FROM `categories` WHERE `name` = 'global';
-- +goose StatementEnd
