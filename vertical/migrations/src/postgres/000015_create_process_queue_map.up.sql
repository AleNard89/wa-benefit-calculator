CREATE TABLE IF NOT EXISTS orchestrator_process_queue_map (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES orgs_companies(id) ON DELETE CASCADE,
    connector_id INTEGER NOT NULL REFERENCES orchestrator_connectors(id) ON DELETE CASCADE,
    process_name VARCHAR(255) NOT NULL,
    queue_name VARCHAR(255) NOT NULL,
    auto_detected BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, connector_id, process_name, queue_name)
);

CREATE INDEX idx_orchestrator_pqm_company_id ON orchestrator_process_queue_map(company_id);

ALTER TABLE orchestrator_process_queue_map ENABLE ROW LEVEL SECURITY;
ALTER TABLE orchestrator_process_queue_map FORCE ROW LEVEL SECURITY;

DO $$ BEGIN
    CREATE POLICY tenant_isolation_orchestrator_process_queue_map ON orchestrator_process_queue_map
        USING (company_id = current_setting('app.current_tenant', true)::integer);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
