-- This migration is not reversible without data loss.
-- Re-add columns if needed, but JSONB data won't be split back.
DROP INDEX IF EXISTS idx_processes_data;
DROP INDEX IF EXISTS idx_processes_results;
ALTER TABLE processes DROP COLUMN IF EXISTS data;
ALTER TABLE processes DROP COLUMN IF EXISTS results;
