BEGIN;

DELETE FROM auth_permissions WHERE app IN ('auth', 'orgs', 'processes');

COMMIT;
