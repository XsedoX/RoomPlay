package com.xsedox.roomplay_backend.domain.rooms.enqueued_songs.vote_statuses;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

/**
 * VoteStatus
 */
public enum VoteStatus {
  UPVOTED("upvoted"),
  DOWNVOTED("downvoted"),
  NOT_VOTED("not_voted");

  private final String value;

  VoteStatus(String value) {
    this.value = value;
  }

  public String getValue() {
    return value;
  }

  public static CustomResult<VoteStatus, DomainError> fromString(String value) {
    for (VoteStatus status : VoteStatus.values()) {
      if (status.getValue().equalsIgnoreCase(value)) {
        return CustomResult.success(status);
      }
    }
    return CustomResult.failure(VoteStatusErrors.INVALID_VOTE_STATUS);
  }

}
