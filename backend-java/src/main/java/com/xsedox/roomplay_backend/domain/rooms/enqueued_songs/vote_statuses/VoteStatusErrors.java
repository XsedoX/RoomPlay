package com.xsedox.roomplay_backend.domain.rooms.enqueued_songs.vote_statuses;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * VoteStatusErrors
 */
public final class VoteStatusErrors {
  public static DomainError INVALID_VOTE_STATUS = new DomainError(
      "VoteStatusErrors.INVALID_VOTE_STATUS", "Invalid vote status provided");
}
