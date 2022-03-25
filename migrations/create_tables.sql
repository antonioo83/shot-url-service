CREATE TABLE IF NOT EXISTS short_url (
    id serial NOT NULL PRIMARY KEY,
    user_code integer NOT NULL,
    correlation_id character varying(50) NOT NULL,
    code character varying(50) NOT NULL,
    original_url character varying(500) NOT NULL,
    short_url character varying(500) NOT NULL,
    CONSTRAINT "shortUrl" UNIQUE (short_url)
);
CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL PRIMARY KEY,
    code integer NOT NULL,
    uid character varying(500) NOT NULL
);