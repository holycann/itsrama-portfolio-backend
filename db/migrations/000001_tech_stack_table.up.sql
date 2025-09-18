-- Create itsrama schema if not exists
CREATE SCHEMA IF NOT EXISTS itsrama;

-- Grant usage and create permissions on schema to service_role
GRANT USAGE, CREATE ON SCHEMA itsrama TO service_role;

-- Create enum type for tech stack categories
CREATE TYPE itsrama.tech_stack_category AS ENUM (
    'Backend',
    'Frontend', 
    'Frameworks',
    'Version Control',
    'Database',
    'DevOps',
    'Tools',
    'CMS & Platforms'
);

CREATE TABLE itsrama.tech_stack (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    category itsrama.tech_stack_category NULL,
    version VARCHAR(50) NULL,
    role VARCHAR(100) NULL,
    is_core_skill BOOLEAN DEFAULT FALSE,
    image_url VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster querying
CREATE INDEX idx_tech_stack_category ON itsrama.tech_stack(category);
CREATE INDEX idx_tech_stack_name ON itsrama.tech_stack(name);

-- Enable Row Level Security
ALTER TABLE itsrama.tech_stack ENABLE ROW LEVEL SECURITY;

-- Grant all permissions on table to service_role
GRANT ALL PRIVILEGES ON TABLE itsrama.tech_stack TO service_role;
GRANT ALL PRIVILEGES ON TYPE itsrama.tech_stack_category TO service_role;

-- Add trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tech_stack_modtime
BEFORE UPDATE ON itsrama.tech_stack
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();
