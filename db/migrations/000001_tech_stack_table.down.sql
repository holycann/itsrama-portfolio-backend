-- Drop trigger
DROP TRIGGER IF EXISTS update_tech_stack_modtime ON itsrama.tech_stack;

-- Drop function
DROP FUNCTION IF EXISTS update_modified_column();

-- Drop indexes
DROP INDEX IF EXISTS itsrama.idx_tech_stack_category;
DROP INDEX IF EXISTS itsrama.idx_tech_stack_name;

-- Drop table
DROP TABLE IF EXISTS itsrama.tech_stack;

-- Drop type
DROP TYPE IF EXISTS itsrama.tech_stack_category;