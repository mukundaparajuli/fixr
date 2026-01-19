CREATE TABLE services_attachment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    service_id UUID REFERENCES services (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    uploaded_by TEXT NOT NULL,
    download_key TEXT NOT NULL,
    file_size BIGINT,
    mime_type TEXT
);

CREATE INDEX idx_services_attachment_service_id ON services_attachment (service_id);

CREATE INDEX idx_services_uploaded_by ON services_attachment (uploaded_by);

CREATE TRIGGER set_services_attachment_updated_at BEFORE
UPDATE ON services_attachment 
FOR EACH ROW 
EXECUTE FUNCTION trigger_set_updated_at();