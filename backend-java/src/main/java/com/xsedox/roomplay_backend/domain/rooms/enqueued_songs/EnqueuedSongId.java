package com.xsedox.roomplay_backend.domain.rooms.enqueued_songs;

import java.util.UUID;

import com.xsedox.roomplay_backend.domain.shared.ConcreteId;

import jakarta.persistence.Embeddable;

/**
 * EnqueuedSongId
 */
@Embeddable
public record EnqueuedSongId(UUID value) implements ConcreteId {
  public static EnqueuedSongId create() {
    return new EnqueuedSongId(UUID.randomUUID());
  }
}
