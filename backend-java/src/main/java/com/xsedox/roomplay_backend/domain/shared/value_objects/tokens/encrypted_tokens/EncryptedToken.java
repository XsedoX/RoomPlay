package com.xsedox.roomplay_backend.domain.shared.value_objects.tokens.encrypted_tokens;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

import jakarta.persistence.Embeddable;

/**
 * EncryptedToken
 */
@Embeddable
public record EncryptedToken(String keyVersion, String value) {
  public EncryptedToken {
    if (value == null || value.isBlank()) {
      throw new IllegalArgumentException(EncryptedTokenErrors.TOKEN_STRING_BLANK_OR_NULL.message());
    }
    if (keyVersion == null || keyVersion.isBlank()) {
      throw new IllegalArgumentException(EncryptedTokenErrors.KEY_VERSION_BLANK_OR_NULL.message());
    }
  }

  public static CustomResult<EncryptedToken, DomainError> create(String keyVersion, String value) {
    if (value == null || value.isBlank()) {
      return CustomResult.failure(EncryptedTokenErrors.TOKEN_STRING_BLANK_OR_NULL);
    }
    if (keyVersion == null || keyVersion.isBlank()) {
      return CustomResult.failure(EncryptedTokenErrors.KEY_VERSION_BLANK_OR_NULL);
    }
    return CustomResult.success(new EncryptedToken(keyVersion, value));
  }
}
