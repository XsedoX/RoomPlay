package infrustructure_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
)

type Seeder struct {
	Queryer contracts.IQueryer
}

func NewSeeder(queryer contracts.IQueryer) *Seeder {
	return &Seeder{
		Queryer: queryer,
	}
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	roomID := uuid.New()
	if err := s.SeedRoom(ctx, roomID); err != nil {
		return err
	}

	userID := uuid.New()
	if err := s.SeedUser(ctx, userID, &roomID); err != nil {
		return err
	}

	songID := uuid.New()
	if err := s.SeedSong(ctx, songID); err != nil {
		return err
	}

	if err := s.SeedEnqueuedSong(ctx, uuid.New(), roomID, songID, userID); err != nil {
		return err
	}

	if err := s.SeedDevice(ctx, uuid.New(), userID); err != nil {
		return err
	}

	// Add more seeds as needed for other tables to be "comprehensive"
	// For now, this covers the main relationships.
	return nil
}

func (s *Seeder) SeedRoom(ctx context.Context, roomID uuid.UUID) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, roomID, "Test Room", []byte("password"), []byte("hash"), time.Now().UTC(), 172800)
	if err != nil {
		return fmt.Errorf("failed to seed room: %w", err)
	}
	return nil
}

func (s *Seeder) SeedUser(ctx context.Context, userID uuid.UUID, roomID *uuid.UUID) error {
	var roomIDVal interface{}
	if roomID != nil {
		roomIDVal = *roomID
	}

	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO users (id, external_id, name, surname, room_id)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, "ext-"+userID.String(), "John", "Doe", roomIDVal)
	if err != nil {
		return fmt.Errorf("failed to seed user: %w", err)
	}
	return nil
}

func (s *Seeder) SeedSong(ctx context.Context, songID uuid.UUID) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO songs (id, external_id, title, author, length_seconds, album_cover_url)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, songID, "ext-song-"+songID.String(), "Test Song", "Test Author", 180, "http://example.com/cover.jpg")
	if err != nil {
		return fmt.Errorf("failed to seed song: %w", err)
	}
	return nil
}

func (s *Seeder) SeedEnqueuedSong(ctx context.Context, id, roomID, songID, userID uuid.UUID) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO enqueued_songs (id, room_id, song_id, added_by, added_at_utc, state, votes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, roomID, songID, userID, time.Now().UTC(), "enqueued", 0)
	if err != nil {
		return fmt.Errorf("failed to seed enqueued song: %w", err)
	}
	return nil
}

func (s *Seeder) SeedDevice(ctx context.Context, deviceID, userID uuid.UUID) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO devices (id, friendly_name, is_host, type, user_id, state, last_logged_in_at_utc)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, deviceID, "Test Device", false, "mobile", userID, "online", time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to seed device: %w", err)
	}
	return nil
}
