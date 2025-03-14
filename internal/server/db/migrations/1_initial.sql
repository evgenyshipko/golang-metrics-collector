-- +goose Up
-- +goose StatementBegin
CREATE TABLE metrics (
                         id SERIAL PRIMARY KEY,
                         name TEXT NOT NULL,
                         type TEXT CHECK (type IN ('gauge', 'counter')) NOT NULL,
                         value_int BIGINT,
                         value_float DOUBLE PRECISION,
                         CHECK (
                             (type = 'counter' AND value_int IS NOT NULL AND value_float IS NULL) OR
                             (type = 'gauge' AND value_float IS NOT NULL AND value_int IS NULL)
                             ),
                         UNIQUE (name, type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE metrics;
-- +goose StatementEnd