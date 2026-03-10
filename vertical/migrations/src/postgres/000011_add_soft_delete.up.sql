ALTER TABLE processes ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;

CREATE INDEX idx_processes_deleted_at ON processes (deleted_at) WHERE deleted_at IS NULL;
