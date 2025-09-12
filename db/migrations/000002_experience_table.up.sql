-- Ensure itsrama schema exists
CREATE SCHEMA IF NOT EXISTS itsrama;

-- Grant usage and create permissions on schema to service_role
GRANT USAGE, CREATE ON SCHEMA itsrama TO service_role;

CREATE TABLE itsrama.experience (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    logo_url TEXT,
    job_type VARCHAR(100),
    start_date DATE NOT NULL,
    end_date DATE,
    location VARCHAR(255),
    arrangement VARCHAR(100),
    work_description TEXT,
    impact TEXT[],
    images_url JSONB[],
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster querying
CREATE INDEX idx_experience_company ON itsrama.experience(company);
CREATE INDEX idx_experience_role ON itsrama.experience(role);
CREATE INDEX idx_experience_start_date ON itsrama.experience(start_date);

-- Enable Row Level Security
ALTER TABLE itsrama.experience ENABLE ROW LEVEL SECURITY;

-- Grant all permissions on table to service_role
GRANT ALL PRIVILEGES ON TABLE itsrama.experience TO service_role;

-- Add trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_experience_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_experience_modtime
BEFORE UPDATE ON itsrama.experience
FOR EACH ROW
EXECUTE FUNCTION update_experience_modified_column();
