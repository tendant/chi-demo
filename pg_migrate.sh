#!/usr/bin/env bash

SOURCE_PATH=migrations/postgres
SOURCE_URL=file://$SOURCE_PATH
PG_USER=demo
PG_PASSWORD=pwd
PG_PORT=5432
PG_NAME=demo_db
PG_DB_URL=postgres://"$PG_USER":"$PG_PASSWORD"@"$PG_HOST":"$PG_PORT"/"$PG_NAME"?sslmode=disable

cd "$(dirname "$0")"

if [ "$1" = "up" ] || [ "$1" = "down" ]
then
  echo "Beginning migration"
  migrate -source $SOURCE_URL -database $PG_DB_URL $1
  echo "Migration complete, beginning dump"
  pg_dump -s -d $PG_NAME -h $PG_HOST -p $PG_PORT -U $PG_USER > $SCHEMA_PATH
  echo "Dump complete, beginning SQLC"
  sqlc generate
  echo "Process complete"
elif test "$1" = "new"
then
  migrate create -dir "$SOURCE_PATH" -seq -digits 6 -ext sql $2
  echo "Migration created"
else
  echo "NOTHING"
fi

