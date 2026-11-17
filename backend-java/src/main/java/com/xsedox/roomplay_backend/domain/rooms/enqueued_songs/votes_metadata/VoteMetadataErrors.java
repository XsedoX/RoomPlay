package com.xsedox.roomplay_backend.domain.rooms.enqueued_songs.votes_metadata;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * VoteMetadataErrors
 */
public final class VoteMetadataErrors {
  public static DomainError VOTED_BY_CANNOT_BE_NULL = new DomainError("VoteMetadataErrors.VOTED_BY_CANNOT_BE_NULL",
      "votedBy must be provided");

  public static DomainError VOTE_STATUS_CANNOT_BE_NULL = new DomainError(
      "VoteMetadataErrors.VOTE_STATUS_CANNOT_BE_NULL",
      "vote status must be provided");
}
