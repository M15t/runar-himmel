DROP SCHEMA IF EXISTS public;
CREATE SCHEMA IF NOT EXISTS main;
ALTER DATABASE maindb SET search_path TO main;
SET search_path to main;

CREATE ROLE dbadmin WITH LOGIN PASSWORD 'DBAdmin123';
GRANT ALL ON DATABASE maindb to dbadmin;
GRANT ALL ON SCHEMA main TO dbadmin;

CREATE ROLE dbreader WITH PASSWORD 'DBReader123';
GRANT CONNECT ON DATABASE maindb TO dbreader;
GRANT USAGE ON SCHEMA main TO dbreader;
ALTER DEFAULT PRIVILEGES FOR ROLE dbadmin IN SCHEMA main GRANT SELECT ON TABLES TO dbreader;

CREATE ROLE dbwriter WITH PASSWORD 'DBWriter123';
GRANT CONNECT ON DATABASE maindb TO dbwriter;
GRANT USAGE ON SCHEMA main TO dbwriter;
ALTER DEFAULT PRIVILEGES FOR ROLE dbadmin IN SCHEMA main GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON TABLES TO dbwriter;
