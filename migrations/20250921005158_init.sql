-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS
  researches (
    id TEXT NOT NULL PRIMARY KEY,
    num SERIAL,
    title TEXT,
    pathology_level REAL, -- от 0 до 1
    created_at TIMESTAMP DEFAULT current_timestamp
  );

CREATE TABLE IF NOT EXISTS
  dicoms (
    id TEXT NOT NULL PRIMARY KEY,
    num SERIAL,
    research_id TEXT,
    file_path TEXT,
    file_size BIGINT,
    assessment TEXT, -- верно, неверно, не оценено
    file_corrupt TEXT, -- ок, полностью битый, частично битый
    probability_of_pathology REAL, -- от 0 до 1
    created_at TIMESTAMP DEFAULT current_timestamp,
    processing_started_at TIMESTAMP,
    processing_finished_at TIMESTAMP,
    FOREIGN KEY (research_id) REFERENCES researches(id) ON DELETE CASCADE
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE researches;
DROP TABLE DICOMS;
-- +goose StatementEnd
