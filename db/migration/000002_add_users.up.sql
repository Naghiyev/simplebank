CREATE TABLE "users"(
    "username" varchar PRIMARY KEY ,
    "hashed_password" varchar NOT NULL ,
    "full_name" varchar NOT NULL ,
    "email" varchar UNIQUE  NOT NULL ,
    "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamp NOT NULL  default 'now()'
);

ALTER TABLE  "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

--CREATE UNIQUE INDEX ON "accounts" ("owner","currency");

ALTER TABLE  "accounts" add constraint "owner_currency_key" unique ("owner","currency")