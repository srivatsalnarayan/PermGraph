-- core system schema
CREATE SCHEMA IF NOT EXISTS zanzibar_core;

-- tenant table
CREATE TABLE zanzibar_core.tenants (
    tenant_id BIGSERIAL PRIMARY KEY,
    tenant_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- revision counter
CREATE TABLE zanzibar_core.revision_counter (
    tenant_id BIGINT PRIMARY KEY,
    current_revision BIGINT
);

-- revision history
CREATE TABLE zanzibar_core.tenant_revision (
    tenant_id BIGINT,
    revision BIGINT,
    valid_from_timestamp TIMESTAMP,
    valid_to_timestamp TIMESTAMP,
    PRIMARY KEY (tenant_id, revision)
);