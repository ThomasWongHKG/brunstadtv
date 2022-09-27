-- +goose Up
/***********************************************************/
/*** SCRIPT AUTHOR: Fredrik Vedvik (fredrik@vedvik.tech) ***/
/***    CREATED ON: 2022-09-27T09:56:41.833Z             ***/
/***********************************************************/

--- BEGIN ALTER TABLE "public"."episodes" ---

ALTER TABLE IF EXISTS "public"."episodes" ADD COLUMN IF NOT EXISTS "parent_id" int4 NULL  ;

COMMENT ON COLUMN "public"."episodes"."parent_id"  IS NULL;

ALTER TABLE IF EXISTS "public"."episodes" ADD CONSTRAINT "episodes_parent_id_foreign" FOREIGN KEY (parent_id) REFERENCES episodes(id);

COMMENT ON CONSTRAINT "episodes_parent_id_foreign" ON "public"."episodes" IS NULL;

--- END ALTER TABLE "public"."episodes" ---

--- BEGIN SYNCHRONIZE TABLE "public"."directus_fields" RECORDS ---

UPDATE "public"."directus_fields" SET "sort" = 9 WHERE "id" = 142;

UPDATE "public"."directus_fields" SET "sort" = 10 WHERE "id" = 123;

UPDATE "public"."directus_fields" SET "sort" = 5, "conditions" = '[{"name":"Don''t require Inheriting","rule":{"_and":[{"type":{"_eq":"inheriting"}}]},"required":false,"options":{"includeSeconds":false,"use24":true}}]' WHERE "id" = 139;

UPDATE "public"."directus_fields" SET "sort" = 7 WHERE "id" = 119;

UPDATE "public"."directus_fields" SET "sort" = 8 WHERE "id" = 130;

INSERT INTO "public"."directus_fields" ("id", "collection", "field", "special", "interface", "options", "display", "display_options", "readonly", "hidden", "sort", "width", "translations", "note", "conditions", "required", "group", "validation", "validation_message")  VALUES (519, 'episodes', 'children', 'o2m', 'list-o2m', '{"filter":{"_and":[{"type":{"_eq":"inheriting"}}]}}', NULL, NULL, false, false, 11, 'full', NULL, NULL, NULL, false, NULL, NULL, NULL);

UPDATE "public"."directus_fields" SET "sort" = 4 WHERE "id" = 128;

UPDATE "public"."directus_fields" SET "sort" = 6 WHERE "id" = 118;

INSERT INTO "public"."directus_fields" ("id", "collection", "field", "special", "interface", "options", "display", "display_options", "readonly", "hidden", "sort", "width", "translations", "note", "conditions", "required", "group", "validation", "validation_message")  VALUES (518, 'episodes', 'parent_id', 'm2o', 'select-dropdown-m2o', NULL, 'related-values', NULL, false, false, 2, 'full', NULL, NULL, '[{"name":"Not required and hidden when not Inheriting","rule":{"_and":[{"type":{"_neq":"inheriting"}}]},"hidden":true,"required":false,"options":{"enableCreate":true,"enableSelect":true}}]', true, 'metadata', NULL, NULL);

UPDATE "public"."directus_fields" SET "options" = '{"choices":[{"text":"Episode","value":"episode"},{"text":"Standalone","value":"standalone"},{"text":"Inheriting","value":"inheriting"}],"icon":"build"}' WHERE "id" = 146;

UPDATE "public"."directus_fields" SET "sort" = 3, "conditions" = '[{"name":"Hide when standalone","rule":{"_and":[{"type":{"_eq":"standalone"}}]},"hidden":true,"required":false,"options":{"enableCreate":true,"enableSelect":true}},{"name":"Don''t require Inheriting","rule":{"_and":[{"type":{"_eq":"inheriting"}}]},"required":false,"options":{"enableCreate":true,"enableSelect":true}}]', "required" = true WHERE "id" = 140;

--- END SYNCHRONIZE TABLE "public"."directus_fields" RECORDS ---

--- BEGIN SYNCHRONIZE TABLE "public"."directus_relations" RECORDS ---

INSERT INTO "public"."directus_relations" ("id", "many_collection", "many_field", "one_collection", "one_field", "one_collection_field", "one_allowed_collections", "junction_field", "sort_field", "one_deselect_action")  VALUES (151, 'episodes', 'parent_id', 'episodes', 'children', NULL, NULL, NULL, NULL, 'nullify');

--- END SYNCHRONIZE TABLE "public"."directus_relations" RECORDS ---
-- +goose Down
/***********************************************************/
/*** SCRIPT AUTHOR: Fredrik Vedvik (fredrik@vedvik.tech) ***/
/***    CREATED ON: 2022-09-27T09:56:43.140Z             ***/
/***********************************************************/

--- BEGIN ALTER TABLE "public"."episodes" ---

ALTER TABLE IF EXISTS "public"."episodes" DROP COLUMN IF EXISTS "parent_id" CASCADE; --WARN: Drop column can occure in data loss!

ALTER TABLE IF EXISTS "public"."episodes" DROP CONSTRAINT IF EXISTS "episodes_parent_id_foreign";

--- END ALTER TABLE "public"."episodes" ---

--- BEGIN SYNCHRONIZE TABLE "public"."directus_fields" RECORDS ---

UPDATE "public"."directus_fields" SET "sort" = 6 WHERE "id" = 119;

UPDATE "public"."directus_fields" SET "sort" = 9 WHERE "id" = 123;

UPDATE "public"."directus_fields" SET "sort" = 7 WHERE "id" = 130;

UPDATE "public"."directus_fields" SET "sort" = 4, "conditions" = NULL WHERE "id" = 139;

UPDATE "public"."directus_fields" SET "sort" = 2, "conditions" = '[{"name":"Hide when standalone","rule":{"_and":[{"type":{"_eq":"standalone"}}]},"hidden":true}]', "required" = false WHERE "id" = 140;

UPDATE "public"."directus_fields" SET "sort" = 8 WHERE "id" = 142;

UPDATE "public"."directus_fields" SET "options" = '{"choices":[{"text":"Episode","value":"episode"},{"text":"Standalone","value":"standalone"}],"icon":"build"}' WHERE "id" = 146;

UPDATE "public"."directus_fields" SET "sort" = 3 WHERE "id" = 128;

UPDATE "public"."directus_fields" SET "sort" = 5 WHERE "id" = 118;

DELETE FROM "public"."directus_fields" WHERE "id" = 519;

DELETE FROM "public"."directus_fields" WHERE "id" = 518;

--- END SYNCHRONIZE TABLE "public"."directus_fields" RECORDS ---

--- BEGIN SYNCHRONIZE TABLE "public"."directus_relations" RECORDS ---

DELETE FROM "public"."directus_relations" WHERE "id" = 151;

--- END SYNCHRONIZE TABLE "public"."directus_relations" RECORDS ---
