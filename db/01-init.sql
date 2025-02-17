CREATE SEQUENCE IF NOT EXISTS account_id;

CREATE TABLE "accounts" (
    "id" int4 NOT NULL DEFAULT nextval('account_id'::regclass),
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "transactions" (
    "id" SERIAL PRIMARY KEY,
    "type" TEXT,
    "status" TEXT,
    "sourcepocketid" INT,
    "destinationpocketid" INT,
    "description" TEXT,
    "amount" FLOAT,
    "currency" TEXT,
    "createdat" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "pockets"(
    "id" SERIAL PRIMARY KEY,
    "name" TEXT,
    "category" TEXT,
    "amount" FLOAT,
    "goal" FLOAT,
    "currency" TEXT,
    "createdat" TIMESTAMP,
    "updatedat" TIMESTAMP,
    "deletedat" TIMESTAMP
);
