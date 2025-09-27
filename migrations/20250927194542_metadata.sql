-- +goose Up
-- +goose StatementBegin
ALTER TABLE researches
ADD COLUMN study_id TEXT,
ADD COLUMN series_id TEXT,
ADD COLUMN files_count INT,
ADD COLUMN inference_error TEXT;

ALTER TABLE researches RENAME COLUMN file_corrupt TO archive_corrupt;
ALTER TABLE researches RENAME COLUMN file_size TO archive_size;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE researches
DROP COLUMN IF EXISTS study_id,
DROP COLUMN IF EXISTS series_id,
DROP COLUMN IF EXISTS files_count,
DROP COLUMN IF EXISTS inference_error;

ALTER TABLE researches RENAME COLUMN archive_corrupt TO file_corrupt;
ALTER TABLE researches RENAME COLUMN archive_size TO file_size;
-- +goose StatementEnd
