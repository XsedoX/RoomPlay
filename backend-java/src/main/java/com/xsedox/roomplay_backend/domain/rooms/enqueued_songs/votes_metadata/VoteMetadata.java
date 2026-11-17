package com.xsedox.roomplay_backend.domain.rooms.enqueued_songs.votes_metadata;

import java.util.Objects;

import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import com.xsedox.roomplay_backend.domain.rooms.enqueued_songs.vote_statuses.VoteStatus;
import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;
import com.xsedox.roomplay_backend.domain.users.UserId;

import jakarta.persistence.Embeddable;

/**
 * VoteMetadata
 */
@Embeddable
public record VoteMetadata(
    UserId votedBy,
    @JdbcTypeCode(SqlTypes.NAMED_ENUM) VoteStatus voteStatus) {
  public VoteMetadata {
    Objects.requireNonNull(votedBy, VoteMetadataErrors.VOTED_BY_CANNOT_BE_NULL.message());
    Objects.requireNonNull(voteStatus, VoteMetadataErrors.VOTE_STATUS_CANNOT_BE_NULL.message());
  }

  public static CustomResult<VoteMetadata, DomainError> create(UserId votedBy, VoteStatus voteStatus) {
    if (votedBy == null) {
      return CustomResult.failure(VoteMetadataErrors.VOTED_BY_CANNOT_BE_NULL);
    }

    if (voteStatus == null) {
      return CustomResult.failure(VoteMetadataErrors.VOTE_STATUS_CANNOT_BE_NULL);
    }

    return CustomResult.success(new VoteMetadata(votedBy, voteStatus));
  }
}
