package com.xsedox.roomplay_backend.domain.shared.value_objects.songs_data;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * SongDataErrors
 */
public final class SongDataErrors {
  public static DomainError EXTERNAL_ID_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "SongDataErrors.EXTERNAL_ID_CANNOT_BE_NULL_OR_BLANK",
      "External id must be provided");

  public static DomainError TITLE_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "SongDataErrors.TITLE_CANNOT_BE_NULL_OR_BLANK",
      "Title must be provided");

  public static DomainError AUTHOR_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "SongDataErrors.AUTHOR_CANNOT_BE_NULL_OR_BLANK",
      "Author must be provided");

  public static DomainError DURATION_SECONDS_MUST_BE_GREATER_THAN_0 = new DomainError(
      "SongDataErrors.DURATION_SECONDS_MUST_BE_GREATER_THAN_0",
      "Duration seconds must be greater than 0");

  public static DomainError ALBUM_COVER_URL_CANNOT_BE_NULL = new DomainError(
      "SongDataErrors.ALBUM_COVER_URL_CANNOT_BE_NULL",
      "Album cover url must be provided");

  public static DomainError MUSIC_PROVIDER_CANNOT_BE_NULL = new DomainError(
      "SongDataErrors.MUSIC_PROVIDER_CANNOT_BE_NULL",
      "Music provider must be provided");

  public static DomainError ISRC_MUST_BE_VALID = new DomainError(
      "SongDataErrors.ISRC_MUST_BE_VALID",
      "ISRC must be valid");

}
