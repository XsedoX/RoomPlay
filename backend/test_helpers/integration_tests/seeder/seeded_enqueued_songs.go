package seeder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/vote_status"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type UsersVotesStruct struct {
	UserId           user_id.UserId
	EnqueuedSongsIds []enqueued_song_id.EnqueuedSongId
	VoteStatus       vote_status.VoteStatus
}

var (
	songsStartedAt = []time.Time{
		time.Date(2001, 11, 22, 12, 5, 0o0, 0o0, time.UTC),
		time.Date(2022, 12, 1, 15, 30, 0o0, 0o0, time.UTC),
		time.Date(2023, 6, 15, 18, 45, 0o0, 0o0, time.UTC),
	}

	isrc0 = "USS1Z2500001"
	isrc1 = "USS1Z2500002"
	isrc2 = "USS1Z2500003"
	isrc3 = "USS1Z2500004"

	enqueuedSongs = []enqueued_song.EnqueuedSong{
		*enqueued_song.HydrateEnqueuedSong(
			enqueued_song_id.NewEnqueuedSongId(),
			*song_data.HydrateSongData(
				faker.URL(),
				faker.Word(),
				faker.Name(),
				faker.URL(),
				349,
				music_provider.YouTube,
				&isrc0,
			),
			time.Date(2001, 11, 22, 12, 0o0, 0o0, 0o0, time.UTC),
			&songsStartedAt[0],
			enqueued_song_state.Played,
			1,
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
				music_provider.YouTube,
				&isrc1,
			),
			time.Date(2001, 10, 22, 12, 0o0, 0o0, 0o0, time.UTC),
			&songsStartedAt[1],
			enqueued_song_state.Playing,
			-2,
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
				music_provider.YouTube,
				&isrc2,
			),
			time.Date(2003, 10, 22, 12, 0o0, 0o0, 0o0, time.UTC),
			nil,
			enqueued_song_state.Enqueued,
			-1,
			userIds[4],
		),
		*enqueued_song.HydrateEnqueuedSong(
			enqueued_song_id.NewEnqueuedSongId(),
			*song_data.HydrateSongData(
				faker.URL(),
				faker.Word(),
				faker.Name(),
				faker.URL(),
				250,
				music_provider.YouTube,
				&isrc3,
			),
			time.Date(2024, 1, 1, 12, 0o0, 0o0, 0o0, time.UTC),
			nil,
			enqueued_song_state.Enqueued,
			-1,
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
				music_provider.YouTube,
				nil,
			),
			time.Date(2023, 5, 10, 12, 0o0, 0o0, 0o0, time.UTC),
			&songsStartedAt[2],
			enqueued_song_state.Played,
			2,
			userIds[3],
		),
	}
	usersVotes = []UsersVotesStruct{
		{
			UserId:           userIds[0],
			EnqueuedSongsIds: []enqueued_song_id.EnqueuedSongId{},
			VoteStatus:       vote_status.NotVoted,
		},
		{
			UserId: userIds[1],
			EnqueuedSongsIds: []enqueued_song_id.EnqueuedSongId{
				enqueuedSongs[4].Id(),
			},
			VoteStatus: vote_status.Upvoted,
		},
		{
			UserId: userIds[1],
			EnqueuedSongsIds: []enqueued_song_id.EnqueuedSongId{
				enqueuedSongs[3].Id(),
			},
			VoteStatus: vote_status.Downvoted,
		},
		{
			UserId: userIds[2],
			EnqueuedSongsIds: []enqueued_song_id.EnqueuedSongId{
				enqueuedSongs[4].Id(),
			},
			VoteStatus: vote_status.Upvoted,
		},
		{
			UserId: userIds[2],
			EnqueuedSongsIds: []enqueued_song_id.EnqueuedSongId{
				enqueuedSongs[3].Id(),
			},
			VoteStatus: vote_status.Downvoted,
		},
		{
			UserId: userIds[3],
			EnqueuedSongsIds: []enqueued_song_id.EnqueuedSongId{
				enqueuedSongs[0].Id(),
			},
			VoteStatus: vote_status.Upvoted,
		},
		{
			UserId: userIds[3],
			EnqueuedSongsIds: []enqueued_song_id.EnqueuedSongId{
				enqueuedSongs[2].Id(),
			},
			VoteStatus: vote_status.Downvoted,
		},
	}
)

func (s *Seeder) seedEnqueuedSong(ctx context.Context, enqueuedSong *enqueued_song.EnqueuedSong, roomId *room_id.RoomId) error {
	var songId uuid.UUID
	errIfSongExists := s.Queryer.QueryRowxContext(ctx,
		`
		INSERT INTO songs (
			id,
			title,
			author,
		  isrc
		)
		VALUES
		(
			$1, $2, $3, $4
		)
		RETURNING id::uuid
		`, uuid.New(),
		enqueuedSong.SongData().Title(),
		enqueuedSong.SongData().Author(),
		enqueuedSong.SongData().Isrc(),
	).Scan(&songId)

	if errors.Is(errIfSongExists, sql.ErrNoRows) {
		err := s.Queryer.QueryRowxContext(ctx,
			`
				SELECT id::uuid
				FROM songs
				WHERE title = $1 AND author = $2 AND isrc = $3;
			`,
			enqueuedSong.SongData().Title(),
			enqueuedSong.SongData().Author(),
			enqueuedSong.SongData().Isrc(),
		).Scan(&songId)
		if err != nil {
			return fmt.Errorf("failed to retrieve existing song id: %w", err)
		}
	} else if errIfSongExists != nil {
		return fmt.Errorf("failed to insert or retrieve song: %w", errIfSongExists)
	}

	_, err := s.Queryer.ExecContext(ctx,
		`
		INSERT INTO songs_external_data
		(
		  song_id,
		  length_seconds,
		  album_cover_url,
			url,
		  music_provider
		)
		VALUES
		(
		  $1, $2, $3, $4, $5
		)
		`,
		songId,
		enqueuedSong.SongData().LengthSeconds(),
		enqueuedSong.SongData().AlbumCoverUrl(),
		enqueuedSong.SongData().Url(),
		enqueuedSong.SongData().MusicProvider().String(),
	)
	if err != nil {
		return fmt.Errorf("failed to seed song external data: %w", err)
	}

	_, err = s.Queryer.ExecContext(ctx,
		`
		INSERT INTO enqueued_songs
		(
			id,
			room_id,
			song_id,
			added_by,
			added_at_utc,
			started_at_utc,
			state
		)
		VALUES
		(
			$1, $2, $3, $4, $5, $6, $7
		)
		`, enqueuedSong.Id(),
		roomId.ToUuid(),
		songId,
		enqueuedSong.AddedBy(),
		enqueuedSong.AddedAtUtc(),
		enqueuedSong.StartedAtUtc(),
		enqueuedSong.State().String(),
	)
	if err != nil {
		return fmt.Errorf("failed to seed enqueued song: %w", err)
	}
	return nil
}

func (s *Seeder) seedUsersVotes(
	ctx context.Context,
	userId user_id.UserId,
	enqueuedSongId enqueued_song_id.EnqueuedSongId,
	voteStatus vote_status.VoteStatus,
) error {
	_, err := s.Queryer.ExecContext(ctx, `
		UPDATE users_votes
		SET vote_status = $1
		WHERE user_id = $2::uuid AND enqueued_song_id = $3::uuid
		`, voteStatus.String(),
		userId.ToUuid(),
		enqueuedSongId.ToUuid(),
	)
	if err != nil {
		return fmt.Errorf("failed to seed user's vote: %w", err)
	}
	return nil
}
