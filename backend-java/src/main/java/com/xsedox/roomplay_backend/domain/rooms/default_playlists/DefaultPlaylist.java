package com.xsedox.roomplay_backend.domain.rooms.default_playlists;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

import jakarta.persistence.Embeddable;

@Embeddable
public record DefaultPlaylist(String externalId, short songsAmount, String title) {
  public DefaultPlaylist {
    if (externalId == null || externalId.isBlank()) {
      throw new IllegalArgumentException("External id cannot be null or blank");
    }
    if (songsAmount < 0) {
      throw new IllegalArgumentException("Songs amount cannot be negative");
    }
    if (title == null || title.isBlank()) {
      throw new IllegalArgumentException("Title cannot be null or blank");
    }
  }

  public static CustomResult<DefaultPlaylist, DomainError> create(String externalId, short songsAmount, String title) {
    if (externalId == null || externalId.isBlank()) {
      return CustomResult.failure(DefaultPlaylistErrors.EXTERNAL_ID_CANNOT_BE_NULL_OR_BLANK);
    }
    if (songsAmount < 0) {
      return CustomResult.failure(DefaultPlaylistErrors.SONGS_AMOUNT_CANNOT_BE_NEGATIVE);
    }
    if (title == null || title.isBlank()) {
      return CustomResult.failure(DefaultPlaylistErrors.TITLE_CANNOT_BE_NULL_OR_BLANK);
    }
    return CustomResult.success(new DefaultPlaylist(externalId, songsAmount, title));
  }
}
