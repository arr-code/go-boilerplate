-- Drop tables in reverse order to avoid foreign key constraint errors
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS user_details;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
