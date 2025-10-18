CREATE TYPE "vote_status" AS ENUM (
  'upvoted',
  'downvoted'
);

CREATE TYPE "song_state" AS ENUM (
  'enqueued',
  'played',
  'playing'
);

CREATE TYPE "user_role" AS ENUM (
  'host',
  'user'
);

CREATE TYPE "device_type" AS ENUM (
  'mobile',
  'computer'
);

CREATE TABLE "songs" (
  "id" uuid PRIMARY KEY,
  "external_id" varchar(256) UNIQUE NOT NULL,
  "title" varchar(256) NOT NULL,
  "author" varchar(256) NOT NULL,
  "mongo_thumbnail_id" uuid NOT NULL,
  "length_seconds" int NOT NULL
);

CREATE TABLE "rooms" (
  "id" uuid PRIMARY KEY,
  "name" varchar(30) NOT NULL,
  "password" varchar(256) NOT NULL,
  "qr_code" varchar(256) NOT NULL,
  "boost_cooldown_seconds" int,
  "created_at_utc" timestamp NOT NULL,
  "lifespan_seconds" int NOT NULL DEFAULT 172800
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "external_id" varchar(256) UNIQUE NOT NULL,
  "email" varchar(256) UNIQUE NOT NULL,
  "room_id" uuid,
  "role" user_role NOT NULL DEFAULT 'user'
);

CREATE TABLE "boosts" (
  "room_id" uuid,
  "user_id" uuid,
  "used_at_utc" time,
  PRIMARY KEY ("room_id", "user_id")
);

CREATE TABLE "enqueued_songs" (
  "id" uuid PRIMARY KEY,
  "room_id" uuid NOT NULL,
  "song_id" uuid NOT NULL,
  "added_by" uuid,
  "added_at_utc" timestamp,
  "played_at_utc" timestamp,
  "state" song_state NOT NULL DEFAULT 'enqueued',
  "votes" int NOT NULL DEFAULT 0
);

CREATE TABLE "rapid_songs" (
  "room_id" uuid,
  "song_id" uuid,
  "to_be_played_at_utc" timestamp NOT NULL,
  PRIMARY KEY ("room_id", "song_id")
);

CREATE TABLE "banned_users" (
  "room_id" uuid,
  "user_id" uuid,
  PRIMARY KEY ("room_id", "user_id")
);

CREATE TABLE "default_playlists" (
  "id" uuid PRIMARY KEY,
  "external_id" varchar(256) UNIQUE NOT NULL,
  "user_id" uuid NOT NULL,
  "song_amount" int NOT NULL,
  "playlist_title" varchar(256) NOT NULL,
  "room_id" uuid NOT NULL
);

CREATE TABLE "devices" (
  "fingerprint" uuid PRIMARY KEY,
  "friendly_name" varchar(30) NOT NULL,
  "isHost" boolean NOT NULL DEFAULT false,
  "type" device_type NOT NULL,
  "user_id" uuid NOT NULL
);

CREATE TABLE "enqueued_songs_users" (
  "enqueued_song_id" uuid,
  "user_id" uuid,
  PRIMARY KEY ("enqueued_song_id", "user_id")
);

COMMENT ON COLUMN "rooms"."lifespan_seconds" IS 'default&max: 48h';

ALTER TABLE "users" ADD CONSTRAINT "user__room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE SET NULL;

ALTER TABLE "boosts" ADD CONSTRAINT "boosts__room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE;

ALTER TABLE "boosts" ADD CONSTRAINT "boosts__user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "enqueued_songs" ADD CONSTRAINT "song_queue__room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE;

ALTER TABLE "enqueued_songs" ADD CONSTRAINT "song_queue__song" FOREIGN KEY ("song_id") REFERENCES "songs" ("id") ON DELETE CASCADE;

ALTER TABLE "enqueued_songs" ADD CONSTRAINT "song_queue__user" FOREIGN KEY ("added_by") REFERENCES "users" ("id") ON DELETE SET NULL;

ALTER TABLE "rapid_songs" ADD CONSTRAINT "rapid_song__room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE;

ALTER TABLE "rapid_songs" ADD CONSTRAINT "rapid_song__song" FOREIGN KEY ("song_id") REFERENCES "songs" ("id") ON DELETE CASCADE;

ALTER TABLE "banned_users" ADD CONSTRAINT "banned_user__room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE;

ALTER TABLE "banned_users" ADD CONSTRAINT "banned_user__user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "default_playlists" ADD CONSTRAINT "default_playlist__room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE;

ALTER TABLE "default_playlists" ADD CONSTRAINT "default_playlist__user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "devices" ADD CONSTRAINT "device__user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "enqueued_songs_users" ADD FOREIGN KEY ("enqueued_song_id") REFERENCES "enqueued_songs" ("id") ON DELETE CASCADE;

ALTER TABLE "enqueued_songs_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
