-- schema owner
CREATE USER practicum WITH password 'practicum';

-- schema user
CREATE USER practicum_ms WITH password 'practicum_ms';

CREATE DATABASE mdb WITH OWNER postgres ENCODING 'UTF8';

CONNECT TO mdb USING postgres;

-- create schema
CREATE SCHEMA practicum AUTHORIZATION practicum;

GRANT USAGE ON SCHEMA practicum TO practicum_ms;

ALTER DEFAULT PRIVILEGES FOR USER practicum IN SCHEMA practicum GRANT SELECT,INSERT,UPDATE,DELETE,TRUNCATE ON TABLES TO practicum;
ALTER DEFAULT PRIVILEGES FOR USER practicum IN SCHEMA practicum GRANT USAGE ON SEQUENCES TO practicum;
ALTER DEFAULT PRIVILEGES FOR USER practicum IN SCHEMA practicum GRANT EXECUTE ON FUNCTIONS TO practicum;
