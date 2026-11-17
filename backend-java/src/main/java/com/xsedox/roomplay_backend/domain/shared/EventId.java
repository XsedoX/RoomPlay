package com.xsedox.roomplay_backend.domain.shared;

import java.util.UUID;

import jakarta.validation.constraints.NotNull;

/**
 * EventId
 */
public record EventId(@NotNull UUID value) implements ConcreteId {
  public EventId {
    if (value == null) {
      throw new IllegalArgumentException("EventId value cannot be null");
    }
  }

  public static EventId create() {
    return new EventId(UUID.randomUUID());
  }

}
