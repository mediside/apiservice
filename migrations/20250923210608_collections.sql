-- +goose Up
-- +goose StatementBegin
ALTER TABLE dicoms RENAME COLUMN research_id TO collection_id;
ALTER TABLE researches RENAME TO collections;
ALTER TABLE dicoms RENAME TO researches;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE researches RENAME TO dicoms;
ALTER TABLE collections RENAME TO researches;
ALTER TABLE dicoms RENAME COLUMN collection_id TO research_id;
-- +goose StatementEnd
