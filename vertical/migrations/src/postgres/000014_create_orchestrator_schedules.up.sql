BEGIN;

CREATE TABLE IF NOT EXISTS orchestrator_schedules (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    connector_id INTEGER NOT NULL,
    external_schedule_id INTEGER,
    name VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    release_name VARCHAR(255),
    package_name VARCHAR(255),
    cron_expression VARCHAR(255),
    cron_summary VARCHAR(500),
    next_occurrence TIMESTAMP WITH TIME ZONE,
    timezone_id VARCHAR(255),
    timezone_iana VARCHAR(255),
    start_strategy INTEGER NOT NULL DEFAULT 0,
    folder_name VARCHAR(255),
    input_arguments JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT orchestrator_schedules_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES orgs_companies (id) ON DELETE CASCADE,
    CONSTRAINT orchestrator_schedules_connector_id_fkey FOREIGN KEY (connector_id)
        REFERENCES orchestrator_connectors (id) ON DELETE CASCADE,
    CONSTRAINT orchestrator_schedules_ext_uq UNIQUE (connector_id, external_schedule_id)
);

CREATE INDEX idx_orchestrator_schedules_company_id ON orchestrator_schedules (company_id);
CREATE INDEX idx_orchestrator_schedules_connector_id ON orchestrator_schedules (connector_id);
CREATE INDEX idx_orchestrator_schedules_enabled ON orchestrator_schedules (enabled);

ALTER TABLE orchestrator_schedules ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_orchestrator_schedules ON orchestrator_schedules
    USING (company_id = current_setting('app.current_tenant', true)::int);
ALTER TABLE orchestrator_schedules FORCE ROW LEVEL SECURITY;

COMMIT;
