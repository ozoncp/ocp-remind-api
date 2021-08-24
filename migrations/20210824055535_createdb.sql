-- +goose Up
-- +goose StatementBegin
CREATE TABLE reminds
(
	id SERIAL PRIMARY KEY,
	remind_id INT,
	user_id INT,
	deadline TIMESTAMP(0) WITH TIME ZONE,
	message TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reminds;
-- +goose StatementEnd
