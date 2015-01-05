-- +migrate Up
ALTER TABLE maps
  ADD sort_by varchar(256) NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE maps
  DROP COLUMN sort_by;
