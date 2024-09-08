CREATE TABLE tenants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at timestamptz NOT NULL DEFAULT 'now()',
    updated_at timestamptz,
    deleted_at timestamptz
);

CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(id) ON DELETE CASCADE NOT NULL,
    amount BIGINT DEFAULT 0 NOT NULL,
    created_at timestamptz NOT NULL DEFAULT 'now()',
    updated_at timestamptz,
    deleted_at timestamptz
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    card_id INT REFERENCES cards(id) ON DELETE CASCADE NOT NULL,
    kind VARCHAR(145) NOT NULL,
    value BIGINT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT 'now()',
    updated_at timestamptz,
    deleted_at timestamptz
);

INSERT INTO tenants (name) VALUES ('Tenant A');
INSERT INTO tenants (name) VALUES ('Tenant B');
INSERT INTO tenants (name) VALUES ('Tenant C');
INSERT INTO tenants (name) VALUES ('Tenant D');
INSERT INTO tenants (name) VALUES ('Tenant E');