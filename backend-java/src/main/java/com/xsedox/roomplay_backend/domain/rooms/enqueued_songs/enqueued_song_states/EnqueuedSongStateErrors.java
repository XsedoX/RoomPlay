package com.xsedox.roomplay_backend.domain.rooms.enqueued_songs.enqueued_song_states;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * EnqueuedSongStateErrors
 */
public final class EnqueuedSongStateErrors {
  public static DomainError INVALID_STATE = new DomainError(
      "EnqueuedSongStateErrors.INVALID_STATE", "Invalid state provided");

}
