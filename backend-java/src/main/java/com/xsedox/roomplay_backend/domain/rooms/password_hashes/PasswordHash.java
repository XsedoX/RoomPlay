package com.xsedox.roomplay_backend.domain.rooms.password_hashes;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

import jakarta.persistence.Embeddable;
import jakarta.validation.constraints.NotNull;

/**
 * PasswordHash
 */
@Embeddable
public record PasswordHash(@NotNull String value) {
  public PasswordHash {
    if (value == null || value.isBlank()) {
      throw new IllegalArgumentException(PasswordHashErrors.PASSWORD_HASH_CANNOT_BE_NULL_OR_BLANK.message());
    }
    if (!value.startsWith("$")) {
      throw new IllegalArgumentException(PasswordHashErrors.PASSWORD_HASH_INVALID_FORMAT.message());
    }
  }

  public static CustomResult<PasswordHash, DomainError> create(String value) {
    if (value == null || value.isBlank()) {
      return CustomResult.failure(PasswordHashErrors.PASSWORD_HASH_CANNOT_BE_NULL_OR_BLANK);
    }
    if (!value.startsWith("$")) {
      return CustomResult.failure(PasswordHashErrors.PASSWORD_HASH_INVALID_FORMAT);
    }
    return CustomResult.success(new PasswordHash(value));
  }
}
