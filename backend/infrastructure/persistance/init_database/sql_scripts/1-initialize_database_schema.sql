DO $$
BEGIN
  -- Check if the 'users' table exists. If it does, we assume everything is initialized.
  IF NOT EXISTS (SELECT FROM pg_tables WHERE tablename = 'users') THEN

CREATE TYPE "music_provider" AS ENUM (
  'youtube',
  'spotify'
);

CREATE TYPE "vote_status" AS ENUM (
  'upvoted',
  'downvoted',
  'not_voted'
);

CREATE TYPE "song_state" AS ENUM (
  'enqueued',
  'played',
  'playing'
);

CREATE TYPE "user_role" AS ENUM (
  'host',
  'member'
);

CREATE TYPE "device_type" AS ENUM (
  'mobile',
  'desktop'
);

CREATE TYPE "device_state" AS ENUM (
  'online',
  'offline'
);

CREATE TABLE "songs_external_data" (
  "song_id" uuid PRIMARY KEY,
  "length_seconds" smallint NOT NULL,
  "album_cover_url" text NOT NULL,
  "url" text UNIQUE NOT NULL,
  "music_provider" music_provider NOT NULL
);

CREATE TABLE "songs" (
  "id" uuid PRIMARY KEY,
  "title" varchar(256) NOT NULL,
  "author" varchar(256) NOT NULL,
  "isrc" varchar(12) UNIQUE
);

CREATE TABLE "rooms" (
  "id" uuid PRIMARY KEY,
  "name" varchar(30) NOT NULL,
  "password" bytea NOT NULL,
  "qr_code_hash" bytea NOT NULL,
  "boost_cooldown_seconds" smallint,
  "created_at_utc" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "lifespan_seconds" int NOT NULL DEFAULT 172800
);

CREATE TABLE "users_votes" (
  "user_id" uuid,
  "enqueued_song_id" uuid,
  "vote_status" vote_status NOT NULL DEFAULT 'not_voted',
  PRIMARY KEY ("user_id", "enqueued_song_id")
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "name" varchar(256) NOT NULL,
  "surname" varchar(256) NOT NULL
);

CREATE TABLE "users_external_credentials" (
  "user_id" uuid PRIMARY KEY,
  "external_id" varchar(256) UNIQUE NOT NULL,
  "access_token" bytea NOT NULL,
  "refresh_token" bytea NOT NULL,
  "music_provider" music_provider NOT NULL,
  "access_token_expires_at_utc" timestamp NOT NULL,
  "refresh_token_expires_at_utc" timestamp NOT NULL,
  "issued_at_utc" timestamp NOT NULL
);

CREATE TABLE "users_internal_credentials" (
  "user_id" uuid,
  "device_id" uuid,
  "refresh_token" bytea UNIQUE NOT NULL,
  "expires_at_utc" timestamp NOT NULL,
  "issued_at_utc" timestamp NOT NULL,
  PRIMARY KEY ("user_id", "device_id")
);

CREATE TABLE "users_room_data" (
  "room_id" uuid,
  "user_id" uuid,
  "boost_used_at_utc" timestamp,
  "role" user_role NOT NULL DEFAULT 'member',
  PRIMARY KEY ("user_id", "room_id")
);

CREATE TABLE "enqueued_songs" (
  "id" uuid PRIMARY KEY,
  "room_id" uuid NOT NULL,
  "song_id" uuid NOT NULL,
  "added_by" uuid NOT NULL,
  "added_at_utc" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "started_at_utc" timestamp,
  "state" song_state NOT NULL DEFAULT 'enqueued'
);

CREATE TABLE "scheduled_songs" (
  "room_id" uuid,
  "song_id" uuid,
  "scheduled_at_utc" timestamp NOT NULL,
  PRIMARY KEY ("room_id", "song_id")
);

CREATE TABLE "banned_users" (
  "room_id" uuid,
  "user_id" uuid,
  PRIMARY KEY ("room_id", "user_id")
);

CREATE TABLE "default_playlists" (
  "external_id" varchar(256) UNIQUE NOT NULL,
  "user_id" uuid NOT NULL,
  "songs_amount" smallint NOT NULL,
  "playlist_title" varchar(256) NOT NULL,
  "room_id" uuid PRIMARY KEY
);

CREATE TABLE "devices" (
  "id" uuid PRIMARY KEY,
  "friendly_name" varchar(30) NOT NULL,
  "is_host" boolean NOT NULL DEFAULT false,
  "type" device_type NOT NULL,
  "user_id" uuid NOT NULL,
  "state" device_state NOT NULL DEFAULT 'online',
  "last_logged_in_at_utc" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE INDEX "songs_url_ix" ON "songs_external_data" ("url");

CREATE INDEX "songs_isrc_ix" ON "songs" ("isrc");

CREATE INDEX "rooms_name_room_spassword_ix" ON "rooms" ("name", "password");

CREATE INDEX "rooms_qr_code_hash_ix" ON "rooms" ("qr_code_hash");

CREATE INDEX "users_external_credentials_external_id_ix" ON "users_external_credentials" ("external_id");

CREATE INDEX "users_internal_credentials_refresh_token_ix" ON "users_internal_credentials" ("refresh_token");

CREATE UNIQUE INDEX "devices_id_user_id_uq" ON "devices" ("id", "user_id");

CREATE INDEX "devices_user_id_ix" ON "devices" ("user_id");

COMMENT ON COLUMN "rooms"."lifespan_seconds" IS 'default&max: 48h';

ALTER TABLE "users_internal_credentials" ADD CONSTRAINT "users__users_refresh_tokens" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "users_internal_credentials" ADD CONSTRAINT "device__users_refresh_tokens" FOREIGN KEY ("device_id") REFERENCES "devices" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "users_external_credentials" ADD CONSTRAINT "users__users_credentials" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "users_room_data" ADD CONSTRAINT "users__users_room_data" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "users_room_data" ADD CONSTRAINT "rooms__users_room_data" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "users_votes" ADD CONSTRAINT "users_votes__user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "users_votes" ADD CONSTRAINT "users_votes__song" FOREIGN KEY ("enqueued_song_id") REFERENCES "enqueued_songs" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "enqueued_songs" ADD CONSTRAINT "song_queue__room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "enqueued_songs" ADD CONSTRAINT "song_queue__song" FOREIGN KEY ("song_id") REFERENCES "songs" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "enqueued_songs" ADD CONSTRAINT "song_queue__user" FOREIGN KEY ("added_by") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "scheduled_songs" ADD CONSTRAINT "scheduled_song__room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "scheduled_songs" ADD CONSTRAINT "scheduled_song__song" FOREIGN KEY ("song_id") REFERENCES "songs" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "banned_users" ADD CONSTRAINT "banned_user__room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "banned_users" ADD CONSTRAINT "banned_user__user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "default_playlists" ADD CONSTRAINT "default_playlist__room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "default_playlists" ADD CONSTRAINT "default_playlist__user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "devices" ADD CONSTRAINT "device__user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "songs_external_data" ADD CONSTRAINT "songs__songs_external_data" FOREIGN KEY ("song_id") REFERENCES "songs" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

  END IF;
END
$$;
