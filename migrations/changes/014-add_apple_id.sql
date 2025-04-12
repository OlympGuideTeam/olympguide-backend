BEGIN;
ALTER TABLE olympguide."user"
    ADD COLUMN apple_id TEXT UNIQUE;
COMMIT;

create index user_apple_id_index
    on olympguide."user" (apple_id);