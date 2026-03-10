ALTER TABLE orgs_companies ADD COLUMN storage_path TEXT NOT NULL DEFAULT '';

-- Set storage_path for existing companies (folder created by API on startup)
UPDATE orgs_companies SET storage_path = '/data/companies/' || LOWER(REGEXP_REPLACE(TRIM(name), '[^a-z0-9]+', '-', 'gi'))
WHERE storage_path = '' OR storage_path IS NULL;
