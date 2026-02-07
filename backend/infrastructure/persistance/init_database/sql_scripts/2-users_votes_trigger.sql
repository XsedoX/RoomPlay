CREATE OR REPLACE FUNCTION populate_users_votes_when_new_song()
RETURNS TRIGGER AS $$
  BEGIN
    INSERT INTO "users_votes" 
    (
      "user_id",
      "enqueued_song_id"
    )
    SELECT 
        urd.user_id,
        NEW.id
    FROM "users_room_data" urd
    WHERE urd.room_id = NEW.room_id;
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

-- 2. Create the Trigger (Idempotent approach)
-- We drop it first to ensure we don't get a duplicate error, then recreate it.
DROP TRIGGER IF EXISTS trigger_auto_create_users_votes ON "enqueued_songs";

CREATE TRIGGER trigger_auto_create_users_votes
AFTER INSERT ON "enqueued_songs"
FOR EACH ROW
EXECUTE FUNCTION populate_users_votes_when_new_song();
