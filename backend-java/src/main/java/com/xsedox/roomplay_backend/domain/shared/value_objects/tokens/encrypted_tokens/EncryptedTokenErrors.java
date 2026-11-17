package com.xsedox.roomplay_backend.domain.shared.value_objects.tokens.encrypted_tokens;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * EncryptedTokenErrors
 */
public final class EncryptedTokenErrors {
  public static DomainError TOKEN_STRING_BLANK_OR_NULL = new DomainError(
      "EncryptedTokenErrors.TOKEN_STRING_BLANK_OR_NULL",
      "Token string must be provided");
  public static DomainError KEY_VERSION_BLANK_OR_NULL = new DomainError(
      "EncryptedTokenErrors.KEY_VERSION_BLANK_OR_NULL",
      "Key version must be provided");

}
