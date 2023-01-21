CREATE SEQUENCE IF NOT EXISTS account_id;

CREATE TABLE "accounts" (
    "id" int4 NOT NULL DEFAULT nextval('account_id'::regclass),
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS transactions(
					id SERIAL PRIMARY KEY,
					`type` TEXT,
					`status` FLOAT,
					souecePocketId TEXT,
					destinationPocketID TEXT,
                    `description` TEXT,
                    amount FLOAT,
                    currency TEXT,
                    createdAt timestamp,
                    );