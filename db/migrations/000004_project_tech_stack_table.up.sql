-- Ensure itsrama schema exists
CREATE SCHEMA IF NOT EXISTS itsrama;

CREATE TABLE itsrama.project_tech_stack (
    project_id UUID NOT NULL REFERENCES itsrama.project(id) ON DELETE CASCADE,
    tech_stack_id UUID NOT NULL REFERENCES itsrama.tech_stack(id) ON DELETE CASCADE,
    
    UNIQUE (project_id, tech_stack_id)
);

-- Create indexes for faster querying
CREATE INDEX idx_project_tech_stack_project ON itsrama.project_tech_stack(project_id);
CREATE INDEX idx_project_tech_stack_tech_stack ON itsrama.project_tech_stack(tech_stack_id);

-- Enable Row Level Security
ALTER TABLE itsrama.project_tech_stack ENABLE ROW LEVEL SECURITY;

-- Grant all permissions on table to service_role
GRANT ALL PRIVILEGES ON TABLE itsrama.project_tech_stack TO service_role;