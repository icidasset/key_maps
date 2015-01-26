-- +migrate Up
ALTER TABLE maps
  DROP COLUMN sort_by;

ALTER TABLE maps
  ADD settings text NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE maps
  DROP COLUMN settings;

ALTER TABLE maps
  ADD sort_by varchar(256) NOT NULL DEFAULT '';
