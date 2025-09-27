-- +goose Up
-- +goose StatementBegin
ALTER TABLE researches ALTER COLUMN archive_corrupt TYPE BOOLEAN USING CASE 
  WHEN archive_corrupt IN ('true', 't', 'yes', 'y', '1') THEN true
  WHEN archive_corrupt IN ('false', 'f', 'no', 'n', '0') THEN false
  ELSE false
END;
ALTER TABLE researches ALTER COLUMN archive_corrupt SET DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE researches ALTER COLUMN archive_corrupt TYPE TEXT;
-- +goose StatementEnd
