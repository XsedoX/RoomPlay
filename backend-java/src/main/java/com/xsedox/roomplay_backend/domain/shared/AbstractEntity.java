package com.xsedox.roomplay_backend.domain.shared;

import jakarta.persistence.AttributeOverride;
import jakarta.persistence.Column;
import jakarta.persistence.EmbeddedId;
import jakarta.persistence.MappedSuperclass;

/**
 * AbstractEntity
 */
@MappedSuperclass
public abstract class AbstractEntity<TId extends ConcreteId> {
  @EmbeddedId
  @AttributeOverride(name = "value", column = @Column(name = "id", nullable = false, updatable = false, unique = true))
  protected TId id;

  protected AbstractEntity() {
  }

  protected AbstractEntity(TId id) {
    this.id = id;
  }

  public TId getId() {
    return id;
  }
}
