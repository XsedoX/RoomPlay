CREATE TYPE music_provider AS ENUM ('spotify', 'youtube');
CREATE TYPE vote_status AS ENUM ('upvoted', 'downvoted', 'not_voted');
CREATE TYPE song_state AS ENUM ('enqueued', 'playing', 'played', 'paused');
CREATE TYPE user_role AS ENUM ('host', 'member');
CREATE TYPE device_type AS ENUM ('desktop', 'mobile');
CREATE TYPE device_state AS ENUM ('online', 'offline');

CREATE TABLE rooms (
  id uuid PRIMARY KEY,
  name varchar(30) NOT NULL,
  password_hash text NOT NULL,
  encrypted_qr_code text NOT NULL,
  boost_cooldown_seconds SMALLINT,
  created_at_utc timestamptz NOT NULL,
  lifespan_seconds integer NOT NULL,
  scheduled_song_external_id varchar(255),
  scheduled_song_title varchar(255),
  scheduled_song_author varchar(255),
  scheduled_song_duration_seconds integer,
  scheduled_song_album_cover_url text,
  scheduled_song_music_provider music_provider,
  scheduled_song_isrc varchar(12),
  scheduled_song_at_utc timestamptz,
  default_playlist_external_id varchar(255),
  default_playlist_songs_amount SMALLINT,
  default_playlist_title varchar(255)
);

CREATE TABLE users (
  id uuid PRIMARY KEY,
  name varchar(255) NOT NULL,
  surname varchar(255) NOT NULL,
  user_role user_role,
  room_id uuid REFERENCES rooms(id) ON DELETE SET NULL,
  boost_used_at_utc timestamptz
);

CREATE TABLE external_credentials (
  user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  external_id varchar(255) UNIQUE NOT NULL,
  access_token text NOT NULL,
  access_token_expires_at_utc timestamptz NOT NULL,
  refresh_token text NOT NULL,
  music_provider music_provider NOT NULL,
  access_token_issued_at_utc timestamptz NOT NULL
);

CREATE TABLE devices (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  friendly_name varchar(30) NOT NULL,
  device_type device_type NOT NULL,
  is_host boolean NOT NULL,
  device_state device_state NOT NULL,
  last_logged_in_utc timestamptz NOT NULL
);

CREATE TABLE enqueued_songs (
  id uuid PRIMARY KEY,
  room_id uuid NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  external_id varchar(255) NOT NULL,
  title varchar(255) NOT NULL,
  author varchar(255) NOT NULL,
  duration_seconds integer NOT NULL,
  album_cover_url text NOT NULL,
  music_provider music_provider NOT NULL,
  isrc varchar(12),
  added_at_utc timestamptz NOT NULL,
  started_at_utc timestamptz,
  song_state song_state NOT NULL,
  added_by_user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE votes_metadata (
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  enqueued_song_id uuid NOT NULL REFERENCES enqueued_songs(id) ON DELETE CASCADE,
  vote_status vote_status NOT NULL,
  PRIMARY KEY (user_id, enqueued_song_id)
);

CREATE TABLE rooms_members (
  room_id uuid NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (room_id, user_id)
);

CREATE TABLE banned_users(
  room_id uuid NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (room_id, user_id)
);

CREATE INDEX idx__rooms__name__password_hash ON rooms(name, password_hash);
CREATE INDEX idx__rooms__encrypted_qr_code ON rooms(encrypted_qr_code);
