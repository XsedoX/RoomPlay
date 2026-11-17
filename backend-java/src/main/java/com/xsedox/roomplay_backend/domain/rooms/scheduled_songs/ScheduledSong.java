package com.xsedox.roomplay_backend.domain.rooms.scheduled_songs;

import java.time.Instant;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;
import com.xsedox.roomplay_backend.domain.shared.value_objects.songs_data.SongData;

import jakarta.persistence.Embeddable;
import jakarta.validation.constraints.NotNull;

/**
 * ScheduledSong
 */
@Embeddable
public record ScheduledSong(SongData songData, @NotNull Instant scheduledAtUtc) {

  public ScheduledSong {
    if (songData == null) {
      throw new IllegalArgumentException("songData cannot be null");
    }
    if (scheduledAtUtc == null) {
      throw new IllegalArgumentException("scheduledTime cannot be null");
    }
  }

  public static CustomResult<ScheduledSong, DomainError> create(SongData songData, Instant scheduledAtUtc,
      Instant currentTime) {
    if (scheduledAtUtc == null) {
      return CustomResult.failure(ScheduledSongErros.SCHEDULED_TIME_CANNOT_BE_NULL);
    }
    if (scheduledAtUtc.isBefore(currentTime)) {
      return CustomResult.failure(ScheduledSongErros.SCHEDULED_TIME_CANNOT_BE_IN_THE_PAST);
    }
    return CustomResult.success(new ScheduledSong(songData, scheduledAtUtc));
  }
}
