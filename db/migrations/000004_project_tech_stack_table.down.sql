-- Drop indexes
DROP INDEX IF EXISTS itsrama.idx_project_tech_stack_project;
DROP INDEX IF EXISTS itsrama.idx_project_tech_stack_tech_stack;

-- Drop project tech stack pivot table
DROP TABLE IF EXISTS itsrama.project_tech_stack;