--liquibase formatted sql

BEGIN;
ALTER TABLE olympguide."user"
    ALTER COLUMN first_name DROP NOT NULL,
    ALTER COLUMN last_name DROP NOT NULL,
    ALTER COLUMN birthday DROP NOT NULL,
    ALTER COLUMN password_hash DROP NOT NULL;

ALTER TABLE olympguide."user"
    ADD COLUMN google_id TEXT UNIQUE,
    ADD COLUMN profile_complete BOOLEAN NOT NULL DEFAULT FALSE;
COMMIT;

create index user_email_index
    on olympguide."user" (email);

create index user_google_id_index
    on olympguide."user" (google_id);