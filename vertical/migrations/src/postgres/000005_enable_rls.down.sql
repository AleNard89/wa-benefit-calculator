BEGIN;

DROP POLICY IF EXISTS tenant_isolation_processes ON processes;
ALTER TABLE processes DISABLE ROW LEVEL SECURITY;

COMMIT;
