CREATE TABLE tenants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO tenants (name) VALUES ('Tenant A');
INSERT INTO tenants (name) VALUES ('Tenant B');
INSERT INTO tenants (name) VALUES ('Tenant C');
INSERT INTO tenants (name) VALUES ('Tenant D');
INSERT INTO tenants (name) VALUES ('Tenant E');

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
    number VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);