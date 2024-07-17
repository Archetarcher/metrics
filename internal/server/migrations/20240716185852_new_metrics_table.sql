-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metrics
(
    id VARCHAR (100) NOT NULL PRIMARY KEY,
    type VARCHAR (10) NOT NULL,
    key VARCHAR (100) NOT NULL,
    delta BIGINT ,
    value DOUBLE PRECISION
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
-- +goose StatementEnd
