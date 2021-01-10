SET client_encoding = 'UTF8';

CREATE TABLE users (
    email text NOT NULL,
    username text NOT NULL,
    name text,
    PRIMARY KEY (email)
);

CREATE TABLE github_oauths (
    email text NOT NULL,
    login text NOT NULL,
    PRIMARY KEY (email),
    FOREIGN KEY (email) REFERENCES users(email) ON DELETE CASCADE
);

CREATE TYPE session_origin AS ENUM ('github', 'facebook', 'apple');

CREATE TABLE sessions (
    access_token text NOT NULL,
    email text NOT NULL,
    refresh_token text,
    origin session_origin NOT NULL,
    PRIMARY KEY (access_token),
    FOREIGN KEY (email) REFERENCES users(email)
)
