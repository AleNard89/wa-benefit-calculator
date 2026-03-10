BEGIN;

CREATE TABLE IF NOT EXISTS auth_permissions (
    id SERIAL PRIMARY KEY,
    app VARCHAR(255) NOT NULL,
    code VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS auth_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_superuser BOOLEAN NOT NULL DEFAULT FALSE,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS auth_users_companies (
    user_id INTEGER NOT NULL,
    company_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, company_id),
    CONSTRAINT auth_users_companies_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES auth_users (id) ON DELETE CASCADE,
    CONSTRAINT auth_users_companies_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES orgs_companies (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS auth_roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(255) NULL
);

CREATE TABLE IF NOT EXISTS auth_users_roles (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    company_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, role_id, company_id),
    CONSTRAINT auth_users_roles_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES auth_users (id) ON DELETE CASCADE,
    CONSTRAINT auth_users_roles_role_id_fkey FOREIGN KEY (role_id)
        REFERENCES auth_roles (id) ON DELETE CASCADE,
    CONSTRAINT auth_users_roles_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES orgs_companies (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS auth_roles_permissions (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    CONSTRAINT auth_roles_permissions_role_id_fkey FOREIGN KEY (role_id)
        REFERENCES auth_roles (id) ON DELETE CASCADE,
    CONSTRAINT auth_roles_permissions_permission_id_fkey FOREIGN KEY (permission_id)
        REFERENCES auth_permissions (id) ON DELETE CASCADE
);

-- Seed admin user (password: Admin123! - bcrypt hash cost 12)
INSERT INTO auth_users (email, password, is_superuser, firstname, lastname, created_at, updated_at)
VALUES ('admin@example.com', '$2a$12$wAP.XV6FOFul1UPdU6O/Tev.lLMCdgag/qcBnwkMQIMGB5K0B4po.', TRUE, 'Admin', 'Admin', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Seed default company
INSERT INTO orgs_companies (id, name, created_at, updated_at)
VALUES (1, 'Demo Company', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Seed Admin role
INSERT INTO auth_roles (id, name, description)
VALUES (1, 'Admin', 'Amministratore con tutti i permessi')
ON CONFLICT DO NOTHING;

-- Assign all permissions to Admin role
INSERT INTO auth_roles_permissions (role_id, permission_id)
SELECT 1, id FROM auth_permissions
ON CONFLICT DO NOTHING;

-- Assign admin user to default company
INSERT INTO auth_users_companies (user_id, company_id)
SELECT u.id, 1 FROM auth_users u WHERE u.email = 'admin@example.com'
ON CONFLICT DO NOTHING;

-- Assign Admin role to admin user for default company
INSERT INTO auth_users_roles (user_id, role_id, company_id)
SELECT u.id, 1, 1 FROM auth_users u WHERE u.email = 'admin@example.com'
ON CONFLICT DO NOTHING;

COMMIT;
