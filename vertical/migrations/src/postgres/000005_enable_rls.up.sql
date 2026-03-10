BEGIN;

ALTER TABLE processes ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_processes ON processes
    USING (company_id = current_setting('app.current_tenant', true)::int);

ALTER TABLE processes FORCE ROW LEVEL SECURITY;

COMMIT;
