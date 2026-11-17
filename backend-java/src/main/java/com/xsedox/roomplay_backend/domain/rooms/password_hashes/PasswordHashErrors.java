package com.xsedox.roomplay_backend.domain.rooms.password_hashes;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * PasswordHashErrors
 */
public final class PasswordHashErrors {
  public static DomainError PASSWORD_HASH_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "PasswordHashErrors.PASSWORD_HASH_CANNOT_BE_NULL_OR_BLANK", "Password hash must be provided");

  public static DomainError PASSWORD_HASH_INVALID_FORMAT = new DomainError(
      "PasswordHashErrors.PASSWORD_HASH_INVALID_FORMAT", "Password hash must be a valid encoded");
}
