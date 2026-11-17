package com.xsedox.roomplay_backend.domain.shared;

import java.time.Instant;

/**
 * DomainEvent
 */
public abstract class AbstractDomainEvent {
  private final Instant occurredAtUtc;
  private final EventId eventId;

  public Instant getOccurredAtUtc() {
    return occurredAtUtc;
  }

  public EventId getEventId() {
    return eventId;
  }

  public AbstractDomainEvent(EventId eventId, Instant occurredAtUtc) {
    this.eventId = eventId;
    this.occurredAtUtc = occurredAtUtc;
  }

}
