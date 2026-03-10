BEGIN;

-- Connectors (UiPath, Python Agent in futuro)
CREATE TABLE IF NOT EXISTS orchestrator_connectors (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('UIPATH', 'PYTHON_AGENT')),
    config JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT orchestrator_connectors_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES orgs_companies (id) ON DELETE CASCADE,
    CONSTRAINT orchestrator_connectors_name_company_uq UNIQUE (company_id, name)
);

CREATE INDEX idx_orchestrator_connectors_company_id ON orchestrator_connectors (company_id);

-- Queue Definitions (code UiPath)
CREATE TABLE IF NOT EXISTS orchestrator_queue_definitions (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    connector_id INTEGER NOT NULL,
    external_definition_id INTEGER,
    name VARCHAR(255) NOT NULL,
    max_retries INTEGER NOT NULL DEFAULT 0,
    folder_name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT orchestrator_queue_definitions_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES orgs_companies (id) ON DELETE CASCADE,
    CONSTRAINT orchestrator_queue_definitions_connector_id_fkey FOREIGN KEY (connector_id)
        REFERENCES orchestrator_connectors (id) ON DELETE CASCADE,
    CONSTRAINT orchestrator_queue_definitions_ext_uq UNIQUE (connector_id, external_definition_id)
);

CREATE INDEX idx_orchestrator_queue_definitions_company_id ON orchestrator_queue_definitions (company_id);
CREATE INDEX idx_orchestrator_queue_definitions_connector_id ON orchestrator_queue_definitions (connector_id);

-- Job Executions (esecuzioni bot)
CREATE TABLE IF NOT EXISTS orchestrator_job_executions (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    connector_id INTEGER NOT NULL,
    external_job_key VARCHAR(255),
    external_job_id INTEGER,
    process_name VARCHAR(255),
    state VARCHAR(50) NOT NULL,
    source_type VARCHAR(50),
    source VARCHAR(255),
    start_time TIMESTAMP WITH TIME ZONE,
    end_time TIMESTAMP WITH TIME ZONE,
    host_machine VARCHAR(255),
    folder_name VARCHAR(255),
    info TEXT,
    details JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT orchestrator_job_executions_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES orgs_companies (id) ON DELETE CASCADE,
    CONSTRAINT orchestrator_job_executions_connector_id_fkey FOREIGN KEY (connector_id)
        REFERENCES orchestrator_connectors (id) ON DELETE CASCADE,
    CONSTRAINT orchestrator_job_executions_ext_uq UNIQUE (connector_id, external_job_key)
);

CREATE INDEX idx_orchestrator_job_executions_company_id ON orchestrator_job_executions (company_id);
CREATE INDEX idx_orchestrator_job_executions_connector_id ON orchestrator_job_executions (connector_id);
CREATE INDEX idx_orchestrator_job_executions_state ON orchestrator_job_executions (state);
CREATE INDEX idx_orchestrator_job_executions_start_time ON orchestrator_job_executions (start_time DESC);

-- Queue Items (elementi nelle code)
CREATE TABLE IF NOT EXISTS orchestrator_queue_items (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    connector_id INTEGER NOT NULL,
    external_item_key VARCHAR(255),
    external_item_id INTEGER,
    queue_definition_id INTEGER,
    queue_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    priority VARCHAR(50) DEFAULT 'Normal',
    reference TEXT,
    processing_exception_type VARCHAR(100),
    error_message TEXT,
    start_processing TIMESTAMP WITH TIME ZONE,
    end_processing TIMESTAMP WITH TIME ZONE,
    retry_number INTEGER NOT NULL DEFAULT 0,
    folder_name VARCHAR(255),
    specific_content JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT orchestrator_queue_items_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES orgs_companies (id) ON DELETE CASCADE,
    CONSTRAINT orchestrator_queue_items_connector_id_fkey FOREIGN KEY (connector_id)
        REFERENCES orchestrator_connectors (id) ON DELETE CASCADE,
    CONSTRAINT orchestrator_queue_items_ext_uq UNIQUE (connector_id, external_item_key)
);

CREATE INDEX idx_orchestrator_queue_items_company_id ON orchestrator_queue_items (company_id);
CREATE INDEX idx_orchestrator_queue_items_connector_id ON orchestrator_queue_items (connector_id);
CREATE INDEX idx_orchestrator_queue_items_status ON orchestrator_queue_items (status);
CREATE INDEX idx_orchestrator_queue_items_queue_name ON orchestrator_queue_items (queue_name);

-- RLS
ALTER TABLE orchestrator_connectors ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_orchestrator_connectors ON orchestrator_connectors
    USING (company_id = current_setting('app.current_tenant', true)::int);
ALTER TABLE orchestrator_connectors FORCE ROW LEVEL SECURITY;

ALTER TABLE orchestrator_queue_definitions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_orchestrator_queue_definitions ON orchestrator_queue_definitions
    USING (company_id = current_setting('app.current_tenant', true)::int);
ALTER TABLE orchestrator_queue_definitions FORCE ROW LEVEL SECURITY;

ALTER TABLE orchestrator_job_executions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_orchestrator_job_executions ON orchestrator_job_executions
    USING (company_id = current_setting('app.current_tenant', true)::int);
ALTER TABLE orchestrator_job_executions FORCE ROW LEVEL SECURITY;

ALTER TABLE orchestrator_queue_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_orchestrator_queue_items ON orchestrator_queue_items
    USING (company_id = current_setting('app.current_tenant', true)::int);
ALTER TABLE orchestrator_queue_items FORCE ROW LEVEL SECURITY;

COMMIT;
