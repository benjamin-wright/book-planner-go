BEGIN;

SAVEPOINT users_restart;

CREATE TABLE users (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "password" varchar(100) NOT NULL,
  "name" varchar(100) NOT NULL UNIQUE
);

RELEASE SAVEPOINT users_restart;

COMMIT;