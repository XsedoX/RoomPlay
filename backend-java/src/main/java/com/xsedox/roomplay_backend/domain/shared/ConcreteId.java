package com.xsedox.roomplay_backend.domain.shared;

import java.util.UUID;

import jakarta.persistence.Embeddable;

/**
 * ConcreteId
 */
@Embeddable
public interface ConcreteId {
  public UUID value();
}
