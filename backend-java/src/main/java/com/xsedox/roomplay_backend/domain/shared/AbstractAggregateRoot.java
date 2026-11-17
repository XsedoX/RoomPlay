package com.xsedox.roomplay_backend.domain.shared;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Objects;

import org.springframework.data.domain.AfterDomainEventPublication;
import org.springframework.data.domain.DomainEvents;

import jakarta.persistence.MappedSuperclass;
import jakarta.persistence.Transient;

/**
 * AbstractAggregateRoot
 */
@MappedSuperclass
public abstract class AbstractAggregateRoot<TId extends ConcreteId> extends AbstractEntity<TId> {
  @Transient
  private final List<AbstractDomainEvent> domainEvents = new ArrayList<>();

  protected AbstractAggregateRoot() {
  }

  protected AbstractAggregateRoot(TId id) {
    super(id);
  }

  protected final void recordEvent(AbstractDomainEvent event) {
    domainEvents.add(Objects.requireNonNull(event));
  }

  @DomainEvents
  public final List<AbstractDomainEvent> getDomainEvents() {
    return Collections.unmodifiableList(domainEvents);
  }

  @AfterDomainEventPublication
  public final void clearDomainEvents() {
    domainEvents.clear();
  }

}
