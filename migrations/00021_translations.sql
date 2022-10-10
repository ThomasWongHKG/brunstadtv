-- +goose Up
/***********************************************************/
/*** SCRIPT AUTHOR: Fredrik Vedvik (fredrik@vedvik.tech) ***/
/***    CREATED ON: 2022-10-10T08:22:35.782Z             ***/
/***********************************************************/

--- BEGIN CREATE OR UPDATE SCHEMA data ---

CREATE SCHEMA IF NOT EXISTS data;

COMMENT ON SCHEMA data  IS NULL;

--- END CREATE OR UPDATE SCHEMA data ---

--- BEGIN CREATE TABLE "data"."translations" ---

CREATE TABLE IF NOT EXISTS "data"."translations" (
	"id" uuid NOT NULL  ,
	"collection" varchar NOT NULL  ,
	"item_id" int4 NOT NULL  ,
	"field" varchar NOT NULL  ,
    "language" varchar NOT NULL  ,
	"value" text NOT NULL  ,
	CONSTRAINT "translations_pk" PRIMARY KEY (id) ,
	CONSTRAINT "translations_unique_key" UNIQUE (collection, item_id, field, "language")
);

CREATE UNIQUE INDEX IF NOT EXISTS translations_unique_key ON data.translations USING btree (collection, item_id, field, "language");


CREATE UNIQUE INDEX IF NOT EXISTS translations_id_uindex ON data.translations USING btree (id);

COMMENT ON COLUMN "data"."translations"."id"  IS NULL;


COMMENT ON COLUMN "data"."translations"."collection"  IS NULL;


COMMENT ON COLUMN "data"."translations"."item_id"  IS NULL;


COMMENT ON COLUMN "data"."translations"."field"  IS NULL;


COMMENT ON COLUMN "data"."translations"."value"  IS NULL;

COMMENT ON CONSTRAINT "translations_pk" ON "data"."translations" IS NULL;


COMMENT ON CONSTRAINT "translations_unique_key" ON "data"."translations" IS NULL;

COMMENT ON INDEX "data"."translations_unique_key"  IS NULL;


COMMENT ON INDEX "data"."translations_id_uindex"  IS NULL;

COMMENT ON TABLE "data"."translations"  IS NULL;

GRANT SELECT ON "data"."translations" TO api, background_worker;
GRANT UPDATE, INSERT, DELETE ON "data"."translations" TO background_worker;

--- END CREATE TABLE "data"."translations" ---
-- +goose Down

DROP TABLE "data"."translations";

DROP SCHEMA "data";
