-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metrics
(
    id VARCHAR (20) NOT NULL PRIMARY KEY,
    type VARCHAR (10) NOT NULL,
    delta INT ,
    value DOUBLE PRECISION
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
-- +goose StatementEnd
