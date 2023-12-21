-- Create database if not exists
CREATE DATABASE IF NOT EXISTS maindb;

-- Switch to the maindb database
USE maindb;

-- Create user
CREATE USER 'dbadmin'@'%' IDENTIFIED BY 'DBAdmin123';

-- Grant privileges to the user on the maindb database
GRANT ALL PRIVILEGES ON maindb.* TO 'dbadmin'@'%';
