-- +migrate Up
CREATE TABLE maps (
  id serial PRIMARY KEY,
  name varchar(256) UNIQUE NOT NULL,
  slug varchar(256) UNIQUE NOT NULL,
  structure text NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  user_id integer NOT NULL
);

CREATE INDEX maps_user_slug_index ON maps(slug)
CREATE INDEX maps_user_id_index ON maps(user_id)

-- +migrate Down
DROP TABLE maps;
