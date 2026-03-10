BEGIN;

-- Auth permissions
INSERT INTO auth_permissions (app, code, description) VALUES ('auth', 'auth:user.read', 'Lettura utenti') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('auth', 'auth:user.create', 'Creazione utenti') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('auth', 'auth:user.update', 'Modifica utenti') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('auth', 'auth:user.delete', 'Eliminazione utenti') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('auth', 'auth:role.read', 'Lettura ruoli') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('auth', 'auth:role.create', 'Creazione ruoli') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('auth', 'auth:role.update', 'Modifica ruoli') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('auth', 'auth:role.delete', 'Eliminazione ruoli') ON CONFLICT DO NOTHING;

-- Orgs permissions
INSERT INTO auth_permissions (app, code, description) VALUES ('orgs', 'orgs:company.read', 'Lettura azienda') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('orgs', 'orgs:company.create', 'Creazione azienda') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('orgs', 'orgs:company.update', 'Modifica azienda') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('orgs', 'orgs:company.delete', 'Eliminazione azienda') ON CONFLICT DO NOTHING;

-- Processes permissions
INSERT INTO auth_permissions (app, code, description) VALUES ('processes', 'processes:process.read', 'Lettura processi') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('processes', 'processes:process.create', 'Creazione processi') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('processes', 'processes:process.update', 'Modifica processi') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('processes', 'processes:process.delete', 'Eliminazione processi') ON CONFLICT DO NOTHING;
INSERT INTO auth_permissions (app, code, description) VALUES ('processes', 'processes:stats.read', 'Lettura statistiche processi') ON CONFLICT DO NOTHING;

COMMIT;
