package com.xsedox.roomplay_backend.domain.rooms;

import java.time.Instant;
import java.util.Collections;
import java.util.HashSet;
import java.util.Optional;
import java.util.Set;

import org.jspecify.annotations.Nullable;

import com.xsedox.roomplay_backend.domain.rooms.default_playlists.DefaultPlaylist;
import com.xsedox.roomplay_backend.domain.rooms.enqueued_songs.EnqueuedSong;
import com.xsedox.roomplay_backend.domain.rooms.password_hashes.PasswordHash;
import com.xsedox.roomplay_backend.domain.rooms.qr_codes.EncryptedQrCode;
import com.xsedox.roomplay_backend.domain.rooms.scheduled_songs.ScheduledSong;
import com.xsedox.roomplay_backend.domain.shared.AbstractAggregateRoot;
import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;
import com.xsedox.roomplay_backend.domain.users.UserId;

import jakarta.persistence.AttributeOverride;
import jakarta.persistence.AttributeOverrides;
import jakarta.persistence.CascadeType;
import jakarta.persistence.CollectionTable;
import jakarta.persistence.Column;
import jakarta.persistence.ElementCollection;
import jakarta.persistence.Embedded;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.OneToMany;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotNull;

/**
 * Room
 */
@Entity
@Table(name = "rooms")
public class Room extends AbstractAggregateRoot<RoomId> {
  protected Room() {
  }

  @NotNull
  @Column(length = 30)
  private String name;
  @Embedded
  @AttributeOverride(name = "value", column = @Column(name = "password_hash"))
  private PasswordHash passwordHash;
  @Embedded
  @AttributeOverrides({
      @AttributeOverride(name = "value", column = @Column(name = "encrypted_qr_code")),
      @AttributeOverride(name = "keyVersion", column = @Column(name = "encrypted_qr_code_key_version"))
  })
  private EncryptedQrCode encryptedQrCode;
  @Nullable
  private Short boostCooldownSeconds;
  @NotNull
  private Instant createdAtUtc;
  private int lifespanSeconds;
  @OneToMany(cascade = CascadeType.ALL, orphanRemoval = true, fetch = FetchType.EAGER)
  @JoinColumn(name = "enqueued_song_id")
  @NotNull
  private Set<EnqueuedSong> enqueuedSongs = new HashSet<>();
  @NotNull
  @ElementCollection(fetch = FetchType.EAGER)
  @CollectionTable(name = "rooms_members", joinColumns = @JoinColumn(name = "room_id"))
  @AttributeOverride(name = "value", column = @Column(name = "user_id"))
  private Set<UserId> members = new HashSet<>();
  @ElementCollection(fetch = FetchType.EAGER)
  @CollectionTable(name = "banned_users", joinColumns = @JoinColumn(name = "room_id"))
  @AttributeOverride(name = "value", column = @Column(name = "user_id"))
  @NotNull
  private Set<UserId> bannedUsers = new HashSet<>();
  @Nullable
  @Embedded
  @AttributeOverrides({
      @AttributeOverride(name = "songData.externalId", column = @Column(name = "scheduled_song_external_id")),
      @AttributeOverride(name = "songData.title", column = @Column(name = "scheduled_song_title")),
      @AttributeOverride(name = "songData.author", column = @Column(name = "scheduled_song_author")),
      @AttributeOverride(name = "songData.durationSeconds", column = @Column(name = "scheduled_song_duration_seconds")),
      @AttributeOverride(name = "songData.albumCoverUrl", column = @Column(name = "scheduled_song_album_cover_url")),
      @AttributeOverride(name = "songData.musicProvider", column = @Column(name = "scheduled_song_music_provider")),
      @AttributeOverride(name = "songData.isrc", column = @Column(name = "scheduled_song_isrc")),
      @AttributeOverride(name = "scheduledAtUtc", column = @Column(name = "scheduled_song_scheduled_at_utc"))
  })
  private ScheduledSong scheduledSong;
  @Nullable
  @Embedded
  @AttributeOverrides({
      @AttributeOverride(name = "externalId", column = @Column(name = "default_playlist_external_id")),
      @AttributeOverride(name = "songsAmount", column = @Column(name = "default_playlist_songs_amount")),
      @AttributeOverride(name = "title", column = @Column(name = "default_playlist_title"))
  })
  private DefaultPlaylist defaultPlaylist;

  public static int DEFAULT_LIFESPAN_SECONDS = 60 * 60 * 24; // NOTE: Default lifespan is 24 hours;

  public String getName() {
    return name;
  }

  public PasswordHash getPasswordHash() {
    return passwordHash;
  }

  public EncryptedQrCode getEncryptedQrCode() {
    return encryptedQrCode;
  }

  public Optional<Short> getBoostCooldownSeconds() {
    return Optional.ofNullable(boostCooldownSeconds);
  }

  public Instant getCreatedAtUtc() {
    return createdAtUtc;
  }

  public int getLifespanSeconds() {
    return lifespanSeconds;
  }

  public Set<EnqueuedSong> getEnqueuedSongs() {
    return Collections.unmodifiableSet(enqueuedSongs);
  }

  public Set<UserId> getMembers() {
    return Collections.unmodifiableSet(members);
  }

  public Set<UserId> getBannedUsers() {
    return Collections.unmodifiableSet(bannedUsers);
  }

  public Optional<ScheduledSong> getScheduledSong() {
    return Optional.ofNullable(scheduledSong);
  }

  public Optional<DefaultPlaylist> getDefaultPlaylist() {
    return Optional.ofNullable(defaultPlaylist);
  }

  private Room(RoomId id, String name, PasswordHash passwordHash, EncryptedQrCode encryptedQrCode,
      @Nullable Short boostCooldownSeconds, Instant createdAtUtc, int lifespanSeconds,
      Set<EnqueuedSong> enqueuedSongs, Set<UserId> members, Set<UserId> bannedUsers,
      @Nullable ScheduledSong scheduledSong, @Nullable DefaultPlaylist defaultPlaylist) {
    super(id);
    this.name = name;
    this.passwordHash = passwordHash;
    this.encryptedQrCode = encryptedQrCode;
    this.boostCooldownSeconds = boostCooldownSeconds;
    this.createdAtUtc = createdAtUtc;
    this.lifespanSeconds = lifespanSeconds;
    this.enqueuedSongs = enqueuedSongs;
    this.members = members;
    this.bannedUsers = bannedUsers;
    this.scheduledSong = scheduledSong;
    this.defaultPlaylist = defaultPlaylist;
  }

  public static CustomResult<Room, DomainError> create(String name, PasswordHash passwordHash,
      EncryptedQrCode encryptedQrCode,
      UserId hostId, Instant currentTimeUtc) {
    Room result = new Room(RoomId.create(), name, passwordHash, encryptedQrCode, null,
        currentTimeUtc, Room.DEFAULT_LIFESPAN_SECONDS, new HashSet<>(), Set.of(hostId),
        new HashSet<>(), null,
        null);

    return CustomResult.success(result);
  }
}
