-- +goose Up 
CREATE TABLE users(
	id UUID,
	created_At TIMESTAMP NOT NULL,
	updated_At TIMESTAMP NOT NULL,
	name TEXT NOT NULL UNIQUE
);
-- +goose Down
DROP TABLE users;
