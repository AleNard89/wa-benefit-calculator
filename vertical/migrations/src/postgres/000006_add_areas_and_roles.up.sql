BEGIN;

-- Areas table (sub-divisions of a company)
CREATE TABLE IF NOT EXISTS orgs_areas (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES orgs_companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, name)
);

CREATE INDEX idx_orgs_areas_company_id ON orgs_areas(company_id);

-- Users can belong to multiple areas
CREATE TABLE IF NOT EXISTS auth_users_areas (
    user_id INTEGER NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    area_id INTEGER NOT NULL REFERENCES orgs_areas(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, area_id)
);

-- Add area_id FK to processes (nullable for migration, replaces text 'area' column)
ALTER TABLE processes ADD COLUMN area_id INTEGER REFERENCES orgs_areas(id) ON DELETE SET NULL;
CREATE INDEX idx_processes_area_id ON processes(area_id);

-- Seed fixed roles: Admin, Contributor, Reader
INSERT INTO auth_roles (id, name, description) VALUES
    (1, 'Admin', 'Amministratore: gestione completa azienda, utenti, ruoli e tutti i processi'),
    (2, 'Contributor', 'Collaboratore: crea e modifica processi nella propria area'),
    (3, 'Reader', 'Lettore: visualizza processi senza metriche finanziarie')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description;

-- Admin role gets all permissions
INSERT INTO auth_roles_permissions (role_id, permission_id)
SELECT 1, id FROM auth_permissions
ON CONFLICT DO NOTHING;

-- Contributor gets process CRUD + stats read
INSERT INTO auth_roles_permissions (role_id, permission_id)
SELECT 2, id FROM auth_permissions WHERE code IN (
    'processes:process.read', 'processes:process.create',
    'processes:process.update', 'processes:process.delete',
    'processes:stats.read'
)
ON CONFLICT DO NOTHING;

-- Reader gets only read permissions
INSERT INTO auth_roles_permissions (role_id, permission_id)
SELECT 3, id FROM auth_permissions WHERE code IN (
    'processes:process.read', 'processes:stats.read'
)
ON CONFLICT DO NOTHING;

-- Seed default areas for company 1
INSERT INTO orgs_areas (company_id, name) VALUES
    (1, 'Finance'),
    (1, 'IT'),
    (1, 'Supply Chain'),
    (1, 'HR'),
    (1, 'Operations')
ON CONFLICT DO NOTHING;

-- Assign admin user to all areas
INSERT INTO auth_users_areas (user_id, area_id)
SELECT u.id, a.id FROM auth_users u, orgs_areas a
WHERE u.email = 'admin@example.com' AND a.company_id = 1
ON CONFLICT DO NOTHING;

COMMIT;
