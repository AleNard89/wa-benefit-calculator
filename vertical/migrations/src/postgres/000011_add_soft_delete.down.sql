DROP INDEX IF EXISTS idx_processes_deleted_at;
ALTER TABLE processes DROP COLUMN IF EXISTS deleted_at;
