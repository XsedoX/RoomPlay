package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/go-faker/faker/v4"
)

var (
	songsStartedAt = []time.Time{
		time.Date(2001, 11, 22, 12, 5, 0o0, 0o0, time.UTC),
		time.Date(2022, 12, 1, 15, 30, 0o0, 0o0, time.UTC),
		time.Date(2023, 6, 15, 18, 45, 0o0, 0o0, time.UTC),
	}
	songs = []enqueued_song.EnqueuedSong{
		*enqueued_song.HydrateEnqueuedSong(
			enqueued_song_id.NewEnqueuedSongId(),
			*song_data.HydrateSongData(
				faker.URL(),
				faker.Word(),
				faker.Name(),
				faker.URL(),
				349,
			),
			time.Date(2001, 11, 22, 12, 0o0, 0o0, 0o0, time.UTC),
			&songsStartedAt[0],
			enqueued_song_state.Played,
			88,
			userIds[0],
		),
		*enqueued_song.HydrateEnqueuedSong(
			enqueued_song_id.NewEnqueuedSongId(),
			*song_data.HydrateSongData(
				faker.URL(),
				faker.Word(),
				faker.Name(),
				faker.URL(),
				349,
			),
			time.Date(2001, 10, 22, 12, 0o0, 0o0, 0o0, time.UTC),
			&songsStartedAt[1],
			enqueued_song_state.Playing,
			8,
			userIds[1],
		),
		*enqueued_song.HydrateEnqueuedSong(
			enqueued_song_id.NewEnqueuedSongId(),
			*song_data.HydrateSongData(
				faker.URL(),
				faker.Word(),
				faker.Name(),
				faker.URL(),
				349,
			),
			time.Date(2003, 10, 22, 12, 0o0, 0o0, 0o0, time.UTC),
			nil,
			enqueued_song_state.Enqueued,
			0,
			userIds[0],
		),
		*enqueued_song.HydrateEnqueuedSong(
			enqueued_song_id.NewEnqueuedSongId(),
			*song_data.HydrateSongData(
				faker.URL(),
				faker.Word(),
				faker.Name(),
				faker.URL(),
				250,
			),
			time.Date(2024, 1, 1, 12, 0o0, 0o0, 0o0, time.UTC),
			nil,
			enqueued_song_state.Enqueued,
			5,
			userIds[2],
		),
		*enqueued_song.HydrateEnqueuedSong(
			enqueued_song_id.NewEnqueuedSongId(),
			*song_data.HydrateSongData(
				faker.URL(),
				faker.Word(),
				faker.Name(),
				faker.URL(),
				180,
			),
			time.Date(2023, 5, 10, 12, 0o0, 0o0, 0o0, time.UTC),
			&songsStartedAt[2],
			enqueued_song_state.Played,
			15,
			userIds[3],
		),
	}
)

func (s *Seeder) seedEnqueuedSong(ctx context.Context, song *enqueued_song.EnqueuedSong) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO songs (id, url, title, author, length_seconds, album_cover_url)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, song.Id(), song.SongData().Url(), song.SongData().Title(), song.SongData().Author(), song.SongData().LengthSeconds(), song.SongData().AlbumCoverUrl())
	if err != nil {
		return fmt.Errorf("failed to seed song: %w", err)
	}
	return nil
}
