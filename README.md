# chi-demo

## Generate sqlc 

    sqlc generate

## Database

### Start Postgres database
    
    # Postgres
    docker run -it --name postgres -e POSTGRES_PASSWORD=pwd -p 5432:5432 postgres:14-alpine
    
    # MS SQL SEVER
    docker run -it  --name mssql -e ACCEPT_EULA='Y' -e MSSQL_SA_PASSWORD='mssql1Pwd' -p 1433:1433 mcr.microsoft.com/azure-sql-edge

### Connect to database using psql

    psql -h localhost -p 5432 -U postgres

### Create postgres database

     CREATE Role demo_project WITH PASSWORD 'pwd';
     grant demo_project to postgres;
     CREATE DATABASE demo_project_db ENCODING 'UTF8' OWNER demo_project;
     GRANT ALL PRIVILEGES ON DATABASE demo_project_db TO demo_project;
     ALTER ROLE demo_project WITH LOGIN;


### Install golang migrate (https://github.com/golang-migrate/migrate)

    brew install golang-migrate 
    
### Create migration

    # If using postgres
    export DATABASE_URL='postgres://demo_project:pwd@localhost:5432/demo_project_db?sslmode=disable'
    migrate create -ext sql -dir migrations/postgres <create_app_user_table>

    # If using MS SQL Server
    export DATABASE_URL="sqlserver://sa:mssql1Pwd@localhost:1433?databaseName=master&integratedSecurity=false&encrypt=false&trustServerCertificate=true"
    migrate create -ext sql -dir migrations/postgres <create_app_user_table>

### Run migration

    migrate -database ${DATABASE_URL} -path migrations/postgres up
    
### Migration Down

    migrate -database ${DATABASE_URL} -path migrations/postgres down

### Migration information table

    schema_migrations

## Fix Dirty DB Migration

    migrate -database ${DATABASE_URL}  -path migrations/postgres force <migration_file_name>

