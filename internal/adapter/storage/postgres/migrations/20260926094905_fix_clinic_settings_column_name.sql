-- +goose Up
-- The original migration named this column lab_name, but every piece of
-- application code (domain.ClinicSettings.ClinicName, the repository's
-- squirrel queries, and the DTOs) has always referred to it as clinic_name —
-- meaning every clinic_settings insert/select/update has been failing with
-- "column clinic_name does not exist" since this table was created.
ALTER TABLE clinic_settings RENAME COLUMN lab_name TO clinic_name;

-- +goose Down
ALTER TABLE clinic_settings RENAME COLUMN clinic_name TO lab_name;
