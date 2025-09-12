-- Drop indexes
DROP INDEX IF EXISTS itsrama.idx_experience_tech_stack_experience;
DROP INDEX IF EXISTS itsrama.idx_experience_tech_stack_tech_stack;

-- Drop experience tech stack pivot table
DROP TABLE IF EXISTS itsrama.experience_tech_stack;