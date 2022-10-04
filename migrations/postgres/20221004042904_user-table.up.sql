BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
   id uuid DEFAULT uuid_generate_v4 (),
   username VARCHAR(20),
   tenant_id VARCHAR (50),
   created_time timestamp NOT NULL default (now() at time zone 'utc')
);

COMMIT;