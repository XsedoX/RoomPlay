package com.xsedox.roomplay_backend.domain.shared.value_objects.tokens;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

public final class TokenErrors {
  public static final DomainError TOKEN_EXPIRATION_TIME_CANNOT_BE_NULL = new DomainError(
      "TokenErrors.TOKEN_EXPIRATION_TIME_CANNOT_BE_NULL",
      "Token expiration time must be provided");
  public static final DomainError TOKEN_VALUE_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "TokenErrors.TOKEN_VALUE_CANNOT_BE_NULL_OR_BLANK",
      "Token value must be provided and cannot be blank");
}
