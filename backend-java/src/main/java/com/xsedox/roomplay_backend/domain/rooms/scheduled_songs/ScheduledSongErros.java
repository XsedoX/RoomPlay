package com.xsedox.roomplay_backend.domain.rooms.scheduled_songs;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

public final class ScheduledSongErros {
  public static DomainError SCHEDULED_TIME_CANNOT_BE_NULL = new DomainError(
      "ScheduledSongErros.SCHEDULED_TIME_CANNOT_BE_NULL", "Scheduled time must be provided");

  public static DomainError SCHEDULED_TIME_CANNOT_BE_IN_THE_PAST = new DomainError(
      "ScheduledSongErros.SCHEDULED_TIME_CANNOT_BE_IN_THE_PAST", "Scheduled time cannot be in the past");

}
