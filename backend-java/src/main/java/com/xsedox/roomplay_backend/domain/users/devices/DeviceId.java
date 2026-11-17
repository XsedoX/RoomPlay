package com.xsedox.roomplay_backend.domain.users.devices;

import java.util.UUID;

import com.xsedox.roomplay_backend.domain.shared.ConcreteId;

import jakarta.persistence.Embeddable;

/**
 * DeviceId
 */
@Embeddable
public record DeviceId(UUID value) implements ConcreteId {
  public static DeviceId create() {
    return new DeviceId(UUID.randomUUID());
  }
}
