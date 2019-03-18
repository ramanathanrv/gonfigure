Goals of the project

Read configuration from Database
 - Mysql
 - Postgres
 - other popular sources

Update the configuration every x minutes
 - x must be configurable by the application

Handle misses gracefully
 - Memoize the response for a short while
 - Ensure the load on the DB is minimal


