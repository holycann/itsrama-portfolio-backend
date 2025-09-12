-- Ensure itsrama schema exists
CREATE SCHEMA IF NOT EXISTS itsrama;

CREATE TABLE itsrama.experience_tech_stack (
    experience_id UUID NOT NULL REFERENCES itsrama.experience(id) ON DELETE CASCADE,
    tech_stack_id UUID NOT NULL REFERENCES itsrama.tech_stack(id) ON DELETE CASCADE,
    
    UNIQUE (experience_id, tech_stack_id)
);

-- Create indexes for faster querying
CREATE INDEX idx_experience_tech_stack_experience ON itsrama.experience_tech_stack(experience_id);
CREATE INDEX idx_experience_tech_stack_tech_stack ON itsrama.experience_tech_stack(tech_stack_id);

-- Enable Row Level Security
ALTER TABLE itsrama.experience_tech_stack ENABLE ROW LEVEL SECURITY;

-- Grant all permissions on table to service_role
GRANT ALL PRIVILEGES ON TABLE itsrama.experience_tech_stack TO service_role;