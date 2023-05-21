BEGIN;

SAVEPOINT books_restart;

CREATE TABLE books (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "name" varchar(100) NOT NULL,
  "summary" varchar(500) NOT NULL,
  "created_time" TIMESTAMP NOT NULL,
  UNIQUE ("user_id", "name")
);

CREATE TABLE characters (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "book_id" UUID NOT NULL,
  "firstname" varchar(50) NOT NULL,
  "familyname" varchar(50),
  CONSTRAINT fk_book
    FOREIGN KEY("book_id")
      REFERENCES books("id")
      ON DELETE CASCADE,
  UNIQUE ("book_id", "firstname", "familyname")
);

CREATE TABLE locations (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "book_id" UUID NOT NULL,
  "name" varchar(50) NOT NULL,
  "description" varchar(500) NOT NULL,
  CONSTRAINT fk_book
    FOREIGN KEY("book_id")
      REFERENCES books("id")
      ON DELETE CASCADE,
  UNIQUE ("book_id", "name")
);

CREATE TABLE chapters (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "book_id" UUID NOT NULL,
  "number" INT NOT NULL,
  "name" VARCHAR(100),
  CONSTRAINT fk_book
    FOREIGN KEY("book_id")
      REFERENCES books("id")
      ON DELETE CASCADE,
  UNIQUE ("book_id", "number")
);

CREATE TABLE scenes (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "chapter_id" UUID NOT NULL,
  "location_id" UUID NOT NULL,
  "number" INT NOT NULL,
  "name" VARCHAR(100) NOT NULL,
  "description" VARCHAR(500) NOT NULL,
  CONSTRAINT fk_chapter
    FOREIGN KEY("chapter_id")
      REFERENCES chapters("id")
      ON DELETE CASCADE,
  UNIQUE ("chapter_id", "number")
);

CREATE TABLE character_scenes (
  "chapter_id" UUID NOT NULL,
  "character_id" UUID NOT NULL,
  CONSTRAINT fk_chapter
    FOREIGN KEY("chapter_id")
      REFERENCES chapters("id")
      ON DELETE CASCADE,
  CONSTRAINT fk_character
    FOREIGN KEY("character_id")
      REFERENCES characters("id")
      ON DELETE CASCADE
);

RELEASE SAVEPOINT books_restart;

COMMIT;