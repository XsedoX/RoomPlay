package com.xsedox.roomplay_backend.domain.rooms.enqueued_songs.enqueued_song_states;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

public enum EnqueuedSongState {
  ENQUEUED("enqueued"),
  PLAYING("playing"),
  PAUSED("paused"),
  PLAYED("played");

  private final String value;

  EnqueuedSongState(String value) {
    this.value = value;
  }

  public String getValue() {
    return value;
  }

  public static CustomResult<EnqueuedSongState, DomainError> fromString(String value) {
    for (EnqueuedSongState state : EnqueuedSongState.values()) {
      if (state.getValue().equalsIgnoreCase(value)) {
        return CustomResult.success(state);
      }
    }
    return CustomResult.failure(EnqueuedSongStateErrors.INVALID_STATE);
  }
}
