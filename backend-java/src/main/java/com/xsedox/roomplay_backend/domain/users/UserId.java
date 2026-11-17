package com.xsedox.roomplay_backend.domain.users;

import java.util.UUID;

import com.xsedox.roomplay_backend.domain.shared.ConcreteId;

import jakarta.persistence.Embeddable;
import jakarta.validation.constraints.NotNull;

@Embeddable
public record UserId(@NotNull UUID value) implements ConcreteId {
  public static UserId create() {
    return new UserId(UUID.randomUUID());
  }
}
