BEGIN;

DROP INDEX IF EXISTS idx_processes_area_id;
ALTER TABLE processes DROP COLUMN IF EXISTS area_id;
DROP TABLE IF EXISTS auth_users_areas;
DROP TABLE IF EXISTS orgs_areas;

DELETE FROM auth_roles WHERE id IN (2, 3);

COMMIT;
