-- use pgcrypto extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS citizen_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    region_id BIGINT NOT NULL,
    district_id BIGINT NOT NULL,
    infrastructure_name VARCHAR(255) NOT NULL,
    sector_id BIGINT NOT NULL,
    description TEXT,
    photo_path VARCHAR(500),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_citizen_reports_region ON citizen_reports(region_id);
CREATE INDEX IF NOT EXISTS idx_citizen_reports_district ON citizen_reports(district_id);
CREATE INDEX IF NOT EXISTS idx_citizen_reports_sector ON citizen_reports(sector_id);

CREATE OR REPLACE TRIGGER update_citizen_reports_updated_at
    BEFORE UPDATE ON citizen_reports
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();