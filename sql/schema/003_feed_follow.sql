-- +goose Up
CREATE TABLE feed_follows (
	id INTEGER PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	user_id UUID NOT NULL,
	feed_id UUID NOT NULL,
	CONSTRAINT fk_feeds FOREIGN KEY (feed_id)
	REFERENCES feeds(id),
	CONSTRAINT fk_users FOREIGN KEY (user_id)
	REFERENCES users(id),
	UNIQUE(feed_id, user_id)
);
-- +goose Down
DROP TABLE feed_follows;
