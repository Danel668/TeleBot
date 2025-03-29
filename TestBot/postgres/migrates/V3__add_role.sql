CREATE TYPE role AS ENUM ('creator', 'admin', 'user');

ALTER TABLE users ADD COLUMN role role;
