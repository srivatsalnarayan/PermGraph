CREATE OR REPLACE FUNCTION zanzibar_core.create_tenant(p_tenant_name TEXT)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    new_tenant_id BIGINT;
    tenant_schema TEXT;
BEGIN

INSERT INTO zanzibar_core.tenants(tenant_name)
VALUES (p_tenant_name)
RETURNING tenant_id INTO new_tenant_id;

INSERT INTO zanzibar_core.revision_counter(tenant_id, current_revision)
VALUES (new_tenant_id, 1);

INSERT INTO zanzibar_core.tenant_revision
VALUES (new_tenant_id, 1, now(), NULL);

tenant_schema := 'tenant_' || new_tenant_id;

EXECUTE format('CREATE SCHEMA %I', tenant_schema);

EXECUTE format(
'CREATE TABLE %I.auth_tuple (
    tuple_id BIGSERIAL PRIMARY KEY,
    object_type TEXT,
    object_id TEXT,
    relation TEXT,
    subject_type TEXT,
    subject_id TEXT,
    valid_from_rev BIGINT,
    valid_to_rev BIGINT
)', tenant_schema);

EXECUTE format(
'CREATE TABLE %I.authorization_model (
    model_id BIGSERIAL PRIMARY KEY,
    model_json JSONB,
    valid_from_rev BIGINT,
    valid_to_rev BIGINT
)', tenant_schema);

RETURN new_tenant_id;

END;
$$;