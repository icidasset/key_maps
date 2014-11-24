-- +migrate Up
CREATE TABLE map_items (
  id serial PRIMARY KEY,
  structure_data text NOT NULL,
  created_at timestamp,
  updated_at timestamp
);

-- +migrate Down
DROP TABLE map_items;
