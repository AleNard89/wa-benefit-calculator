DROP INDEX IF EXISTS idx_orgs_companies_parent_id;
ALTER TABLE orgs_companies DROP COLUMN IF EXISTS parent_id;
