# BG Dictionary API Server

A simple Go server to fetch and serve words from the dictionary database.

## Dev

All database queries are generated with SQLC.

The schema can be altered in `db/schema.sql` and new queries can be added in `db/query.sql`.

The sqlc generated code can then be updated using the following command in the `db` directory:

`docker run --rm -v .:/src -w /src sqlc/sqlc generate`
