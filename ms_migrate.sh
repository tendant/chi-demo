#!/usr/bin/env bash

SOURCE_PATH="migrations/mssql"
SOURCE_URL=file://"$SOURCE_PATH"
MSSQL_DB_URL="sqlserver://sa:mssql1Ipw@localhost:1433?databaseName=TorpagoApp&integratedSecurity=false&encrypt=false&trustServerCertificate=true"

cd "$(dirname "$0")"

if [ "$1" = "up" ] || [ "$1" = "down" ]
then
  echo "Beginning migration"
  migrate -source $SOURCE_URL -database $MSSQL_DB_URL $1
  echo "Migration complete"
elif test "$1" = "new"
then
  migrate create -dir "$SOURCE_PATH" -seq -digits 6 -ext sql $2
  echo "Migration created"
else
  echo "NOTHING"
fi

