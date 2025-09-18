-- Drop trigger
DROP TRIGGER IF EXISTS update_project_modtime ON itsrama.project;

-- Drop function
DROP FUNCTION IF EXISTS update_project_modified_column();

-- Drop indexes
DROP INDEX IF EXISTS itsrama.idx_project_slug;
DROP INDEX IF EXISTS itsrama.idx_project_title;
DROP INDEX IF EXISTS itsrama.idx_project_category;
DROP INDEX IF EXISTS itsrama.idx_project_development_status;
DROP INDEX IF EXISTS itsrama.idx_project_progress_status;

-- Drop table
DROP TABLE IF EXISTS itsrama.project;

-- Drop enum types
DROP TYPE IF EXISTS itsrama.development_status;
DROP TYPE IF EXISTS itsrama.progress_status;
DROP TYPE IF EXISTS itsrama.project_category;