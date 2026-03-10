BEGIN;

CREATE TABLE IF NOT EXISTS processes (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES orgs_companies(id) ON DELETE CASCADE,

    -- Info generali
    process_name VARCHAR(255) NOT NULL,
    process_description TEXT,
    proposer VARCHAR(255) NOT NULL,
    area VARCHAR(255) NOT NULL,
    responsible_manager VARCHAR(255) NOT NULL,
    department VARCHAR(255),

    -- Caratteristiche processo
    systems_involved INTEGER NOT NULL DEFAULT 1,
    process_type VARCHAR(100) NOT NULL,
    periodicity VARCHAR(50) NOT NULL,
    frequent_changes BOOLEAN NOT NULL DEFAULT FALSE,
    technology VARCHAR(100) NOT NULL,

    -- Costi
    implementation_cost DECIMAL(12,2) NOT NULL,
    training_cost DECIMAL(12,2) NOT NULL DEFAULT 0,
    maintenance_cost DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Parametri operativi
    hourly_cost DECIMAL(8,2) NOT NULL,
    time_per_activity INTEGER NOT NULL,
    activities_per_day INTEGER NOT NULL,
    working_days_per_year INTEGER NOT NULL DEFAULT 220,

    -- Errori
    current_error_rate DECIMAL(5,2) NOT NULL,
    post_error_rate DECIMAL(5,2) NOT NULL,
    error_cost DECIMAL(10,2) NOT NULL,

    -- Produttivita
    productivity_factor DECIMAL(4,1) NOT NULL DEFAULT 2.0,
    time_reduction_factor INTEGER NOT NULL DEFAULT 50,

    -- Punteggi impatto (1-5)
    data_quality_score INTEGER NOT NULL DEFAULT 3,
    audit_score INTEGER NOT NULL DEFAULT 3,
    customer_experience_score INTEGER NOT NULL DEFAULT 3,
    error_reduction_score INTEGER NOT NULL DEFAULT 3,
    standardization_score INTEGER NOT NULL DEFAULT 3,
    scalability_score INTEGER NOT NULL DEFAULT 3,

    -- Risultati calcolati
    operational_savings DECIMAL(12,2),
    error_reduction_savings DECIMAL(12,2),
    productivity_benefit DECIMAL(12,2),
    annual_savings DECIMAL(12,2),
    roi DECIMAL(8,2),
    break_even_months INTEGER,
    hours_saved_monthly DECIMAL(8,2),
    hours_saved_annually DECIMAL(8,2),
    impact_score DECIMAL(4,2),

    -- Stato e metadata
    status VARCHAR(50) NOT NULL DEFAULT 'To Valuate',
    created_by INTEGER REFERENCES auth_users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_processes_company_id ON processes(company_id);
CREATE INDEX idx_processes_status ON processes(status);
CREATE INDEX idx_processes_created_at ON processes(created_at DESC);

COMMIT;
