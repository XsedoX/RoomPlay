package com.xsedox.roomplay_backend.domain.rooms;

import java.util.UUID;

import com.sun.istack.NotNull;
import com.xsedox.roomplay_backend.domain.shared.ConcreteId;

import jakarta.persistence.Embeddable;

/**
 * RoomId
 */
@Embeddable
public record RoomId(@NotNull UUID value) implements ConcreteId {
  public static RoomId create() {
    return new RoomId(UUID.randomUUID());
  }
}
