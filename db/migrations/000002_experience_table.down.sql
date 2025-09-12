-- Drop trigger
DROP TRIGGER IF EXISTS update_experience_modtime ON itsrama.experience;

-- Drop function
DROP FUNCTION IF EXISTS update_experience_modified_column();

-- Drop indexes
DROP INDEX IF EXISTS itsrama.idx_experience_company;
DROP INDEX IF EXISTS itsrama.idx_experience_role;
DROP INDEX IF EXISTS itsrama.idx_experience_start_date;

-- Drop table
DROP TABLE IF EXISTS itsrama.experience;