-- +goose Up
-- +goose StatementBegin
INSERT INTO tenants (name) VALUES ('Tenant A');
INSERT INTO tenants (name) VALUES ('Tenant B');
INSERT INTO tenants (name) VALUES ('Tenant C');
INSERT INTO tenants (name) VALUES ('Tenant D');
INSERT INTO tenants (name) VALUES ('Tenant E');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tenants WHERE name IN ('Tenant A', 'Tenant B', 'Tenant C', 'Tenant D', 'Tenant E');
-- +goose StatementEnd
