package com.xsedox.roomplay_backend.domain.rooms.enqueued_songs;

import java.time.Instant;
import java.util.Collections;
import java.util.HashSet;
import java.util.Optional;
import java.util.Set;

import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;
import org.jspecify.annotations.Nullable;

import com.xsedox.roomplay_backend.domain.rooms.enqueued_songs.enqueued_song_states.EnqueuedSongState;
import com.xsedox.roomplay_backend.domain.rooms.enqueued_songs.votes_metadata.VoteMetadata;
import com.xsedox.roomplay_backend.domain.shared.AbstractAggregateRoot;
import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;
import com.xsedox.roomplay_backend.domain.shared.value_objects.songs_data.SongData;
import com.xsedox.roomplay_backend.domain.users.UserId;

import jakarta.persistence.AttributeOverride;
import jakarta.persistence.AttributeOverrides;
import jakarta.persistence.CollectionTable;
import jakarta.persistence.Column;
import jakarta.persistence.ElementCollection;
import jakarta.persistence.Embedded;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotNull;

/**
 * EnqueuedSong
 */
@Entity
@Table(name = "enqueued_songs")
public class EnqueuedSong extends AbstractAggregateRoot<EnqueuedSongId> {
  @Embedded
  private SongData songData;
  @NotNull
  @Column(updatable = false)
  private Instant addedAtUtc;
  @Nullable
  private Instant startedAtUtc;
  @JdbcTypeCode(SqlTypes.NAMED_ENUM)
  @NotNull
  private EnqueuedSongState state;
  @Embedded
  @ElementCollection(fetch = FetchType.EAGER)
  @NotNull
  @CollectionTable(name = "users_votes", joinColumns = @JoinColumn(name = "enqueued_song_id"))
  @AttributeOverrides({
      @AttributeOverride(name = "userId.value", column = @Column(name = "user_id")),
      @AttributeOverride(name = "vote_status", column = @Column(name = "vote_status"))
  })
  private Set<VoteMetadata> votes = new HashSet<>();
  @NotNull
  @Embedded
  @NotNull
  @AttributeOverride(name = "value", column = @Column(name = "added_by_user_id"))
  private UserId addedByUserId;

  public SongData getSongData() {
    return songData;
  }

  public Instant getAddedAtUtc() {
    return addedAtUtc;
  }

  public Optional<Instant> getStartedAtUtc() {
    return Optional.ofNullable(startedAtUtc);
  }

  public EnqueuedSongState getState() {
    return state;
  }

  public Set<VoteMetadata> getVotes() {
    return Collections.unmodifiableSet(votes);
  }

  public UserId getAddedByUserId() {
    return addedByUserId;
  }

  private EnqueuedSong(EnqueuedSongId id, SongData songData, Instant addedAtUtc, @Nullable Instant startedAtUtc,
      EnqueuedSongState state,
      Set<VoteMetadata> votes, UserId addedByUserId) {
    super(id);
    this.songData = songData;
    this.addedAtUtc = addedAtUtc;
    this.startedAtUtc = startedAtUtc;
    this.state = state;
    this.votes = votes;
    this.addedByUserId = addedByUserId;
  }

  public static CustomResult<EnqueuedSong, DomainError> create(SongData songData,
      UserId addedByUserId, Instant currentTimeUtc) {
    EnqueuedSong result = new EnqueuedSong(EnqueuedSongId.create(), songData, currentTimeUtc, null,
        EnqueuedSongState.ENQUEUED, new HashSet<>(), addedByUserId);
    return CustomResult.success(result);
  }
}
