BEGIN;

CREATE TABLE IF NOT EXISTS student
(
    id varchar(255) not null,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    national_id varchar(255) not null
);

CREATE UNIQUE INDEX IF NOT EXISTS id_unique
    on student (id);

CREATE INDEX IF NOT EXISTS name_idx
    on student (first_name);

COMMIT;