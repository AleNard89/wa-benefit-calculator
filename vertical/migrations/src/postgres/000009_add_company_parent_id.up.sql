ALTER TABLE orgs_companies ADD COLUMN parent_id INTEGER REFERENCES orgs_companies(id) ON DELETE SET NULL;
CREATE INDEX idx_orgs_companies_parent_id ON orgs_companies(parent_id);
