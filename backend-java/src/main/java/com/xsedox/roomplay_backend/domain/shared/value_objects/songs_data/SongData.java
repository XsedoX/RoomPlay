package com.xsedox.roomplay_backend.domain.shared.value_objects.songs_data;

import java.net.URI;
import java.util.Objects;
import java.util.Optional;
import java.util.regex.Pattern;

import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;
import org.jspecify.annotations.Nullable;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;
import com.xsedox.roomplay_backend.domain.shared.value_objects.music_providers.MusicProvider;

import jakarta.persistence.Embeddable;
import jakarta.validation.constraints.NotNull;

@Embeddable
public class SongData {
  protected SongData() {
  }

  @NotNull
  private String externalId;
  @NotNull
  private String title;
  @NotNull
  private String author;
  @NotNull
  private int durationSeconds;
  @NotNull
  private URI albumCoverUrl;
  @JdbcTypeCode(SqlTypes.NAMED_ENUM)
  @NotNull
  private MusicProvider musicProvider;
  private @Nullable String isrc;

  public String externalId() {
    return externalId;
  }

  public String title() {
    return title;
  }

  public String author() {
    return author;
  }

  public int durationSeconds() {
    return durationSeconds;
  }

  public URI albumCoverUrl() {
    return albumCoverUrl;
  }

  public MusicProvider musicProvider() {
    return musicProvider;
  }

  public Optional<String> isrc() {
    return Optional.ofNullable(isrc);
  }

  private SongData(
      String externalId,
      String title,
      String author,
      int durationSeconds,
      URI albumCoverUrl,
      MusicProvider musicProvider,
      @Nullable String isrc) {
    if (externalId == null || externalId.isBlank()) {
      throw new IllegalArgumentException(SongDataErrors.EXTERNAL_ID_CANNOT_BE_NULL_OR_BLANK.message());
    }

    if (title == null || title.isBlank()) {
      throw new IllegalArgumentException(SongDataErrors.TITLE_CANNOT_BE_NULL_OR_BLANK.message());
    }

    if (author == null || author.isBlank()) {
      throw new IllegalArgumentException(SongDataErrors.AUTHOR_CANNOT_BE_NULL_OR_BLANK.message());
    }

    if (durationSeconds <= 0) {
      throw new IllegalArgumentException(SongDataErrors.DURATION_SECONDS_MUST_BE_GREATER_THAN_0.message());
    }

    if (albumCoverUrl == null) {
      throw new IllegalArgumentException(SongDataErrors.ALBUM_COVER_URL_CANNOT_BE_NULL.message());
    }

    if (musicProvider == null) {
      throw new IllegalArgumentException(SongDataErrors.MUSIC_PROVIDER_CANNOT_BE_NULL.message());
    }
  }

  @Override
  public int hashCode() {
    return Objects.hash(externalId, title, author, durationSeconds, albumCoverUrl, musicProvider, isrc);
  }

  @Override
  public boolean equals(Object o) {
    if (this == o)
      return true;

    if (!(o instanceof SongData other))
      return false;
    return externalId.equals(other.externalId)
        && title.equals(other.title)
        && author.equals(other.author)
        && durationSeconds == other.durationSeconds
        && albumCoverUrl.equals(other.albumCoverUrl)
        && musicProvider == other.musicProvider
        && isrc == other.isrc;
  }

  public static CustomResult<SongData, DomainError> create(
      String externalId,
      String title,
      String author,
      int durationSeconds,
      URI albumCoverUrl,
      MusicProvider musicProvider,
      @Nullable String isrc) {

    if (externalId == null || externalId.isBlank()) {
      return CustomResult.failure(SongDataErrors.EXTERNAL_ID_CANNOT_BE_NULL_OR_BLANK);
    }

    if (title == null || title.isBlank()) {
      return CustomResult.failure(SongDataErrors.TITLE_CANNOT_BE_NULL_OR_BLANK);
    }

    if (author == null || author.isBlank()) {
      return CustomResult.failure(SongDataErrors.AUTHOR_CANNOT_BE_NULL_OR_BLANK);
    }

    if (durationSeconds <= 0) {
      return CustomResult.failure(SongDataErrors.DURATION_SECONDS_MUST_BE_GREATER_THAN_0);
    }

    if (albumCoverUrl == null) {
      return CustomResult.failure(SongDataErrors.ALBUM_COVER_URL_CANNOT_BE_NULL);
    }

    if (musicProvider == null) {
      return CustomResult.failure(SongDataErrors.MUSIC_PROVIDER_CANNOT_BE_NULL);
    }

    if (isrc != null) {
      Pattern ISRC_PATTERN = Pattern.compile("^[A-Z]{2}[A-Z0-9]{3}[0-9]{2}[0-9]{5}$",
          Pattern.CASE_INSENSITIVE);

      if (!ISRC_PATTERN.matcher(isrc).matches()) {
        return CustomResult.failure(SongDataErrors.ISRC_MUST_BE_VALID);
      }
    }

    return CustomResult.success(new SongData(
        externalId,
        title,
        author,
        durationSeconds,
        albumCoverUrl,
        musicProvider,
        isrc));

  }

}
