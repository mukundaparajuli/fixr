CREATE TABLE services_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    description TEXT,
    color TEXT DEFAULT '#FFFFFF',
    user_id TEXT NOT NULL
);

CREATE INDEX idx_services_categories_user_id ON services_categories (user_id);

CREATE UNIQUE INDEX services_categories_unique_name ON services_categories (user_id, name);

CREATE TRIGGER set_updated_at_services_categories
BEFORE UPDATE ON services_categories
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    rate NUMERIC(5, 2) NOT NULL DEFAULT 0.00,
    method TEXT NOT NULL DEFAULT 'hourly',
    category_id UUID REFERENCES services_categories (id) ON DELETE SET NULL,
    parent_service_id UUID REFERENCES services (id),
    user_id TEXT NOT NULL,
    sort_order SERIAL
);

CREATE INDEX idx_services_user_id ON services (user_id);

CREATE INDEX idx_services_category_id ON services (category_id);

CREATE INDEX idx_services_parent_service_id ON services (parent_service_id);

CREATE INDEX idx_services_status ON services (status);

CREATE UNIQUE INDEX services_unique_name ON services (user_id, name);

CREATE TRIGGER set_services_updated_at
BEFORE UPDATE ON services
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

-- constrains
ALTER TABLE services
ADD CONSTRAINT no_self_parent CHECK (id != parent_service_id);

CREATE INDEX idx_services_hierarchy ON services (parent_service_id, id);