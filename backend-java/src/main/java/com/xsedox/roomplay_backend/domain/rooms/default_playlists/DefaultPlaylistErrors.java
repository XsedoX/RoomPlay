package com.xsedox.roomplay_backend.domain.rooms.default_playlists;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * DefaultPlaylistErrors
 */
public final class DefaultPlaylistErrors {
  public static final DomainError EXTERNAL_ID_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "DefaultPlaylistErrors.EXTERNAL_ID_CANNOT_BE_NULL_OR_BLANK",
      "External id must be provided");

  public static final DomainError SONGS_AMOUNT_CANNOT_BE_NEGATIVE = new DomainError(
      "DefaultPlaylistErrors.SONGS_AMOUNT_CANNOT_BE_NEGATIVE",
      "Songs amount must be greater than 0");

  public static final DomainError TITLE_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "DefaultPlaylistErrors.TITLE_CANNOT_BE_NULL_OR_BLANK",
      "Title must be provided");

}
