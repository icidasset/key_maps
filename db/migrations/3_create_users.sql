-- +migrate Up
CREATE TABLE users (
  id serial PRIMARY KEY,
  email varchar(256) UNIQUE NOT NULL,
  encrypted_password text NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL
);

CREATE INDEX users_email_index ON users(email)

-- +migrate Down
DROP TABLE users;
