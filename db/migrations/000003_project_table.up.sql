-- Ensure itsrama schema exists
CREATE SCHEMA IF NOT EXISTS itsrama;

-- Grant usage and create permissions on schema to service_role
GRANT USAGE, CREATE ON SCHEMA itsrama TO service_role;

-- Enum types for development, progress, and project category status
CREATE TYPE itsrama.development_status AS ENUM ('Alpha', 'Beta', 'MVP');
CREATE TYPE itsrama.progress_status AS ENUM ('In Progress', 'In Revision', 'On Hold', 'Completed');
CREATE TYPE itsrama.project_category AS ENUM ('Web Development', 'API Development', 'Bot Development', 'Mobile App', 'Desktop App', 'UI/UX Design', 'Other');

-- Grant permissions on enum types to service_role
GRANT USAGE ON TYPE itsrama.development_status TO service_role;
GRANT USAGE ON TYPE itsrama.progress_status TO service_role;
GRANT USAGE ON TYPE itsrama.project_category TO service_role;

CREATE TABLE itsrama.project (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255),
    description TEXT NOT NULL,
    my_role TEXT[] NOT NULL,
    category itsrama.project_category,
    
    github_url TEXT,
    web_url TEXT,
    
    images JSONB[],
    features TEXT[],
    
    development_status itsrama.development_status,
    progress_status itsrama.progress_status,
    progress_percentage INTEGER CHECK (progress_percentage BETWEEN 0 AND 100),
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster querying
CREATE INDEX idx_project_slug ON itsrama.project(slug);
CREATE INDEX idx_project_title ON itsrama.project(title);
CREATE INDEX idx_project_category ON itsrama.project(category);
CREATE INDEX idx_project_development_status ON itsrama.project(development_status);
CREATE INDEX idx_project_progress_status ON itsrama.project(progress_status);

-- Enable Row Level Security
ALTER TABLE itsrama.project ENABLE ROW LEVEL SECURITY;

-- Grant all permissions on table to service_role
GRANT ALL PRIVILEGES ON TABLE itsrama.project TO service_role;

-- Add trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_project_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_project_modtime
BEFORE UPDATE ON itsrama.project
FOR EACH ROW
EXECUTE FUNCTION update_project_modified_column();
