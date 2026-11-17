package com.xsedox.roomplay_backend.domain.shared.value_objects.tokens;

import java.time.Instant;
import java.util.Objects;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;
import com.xsedox.roomplay_backend.domain.shared.value_objects.tokens.encrypted_tokens.EncryptedToken;

import jakarta.persistence.Embeddable;
import jakarta.validation.constraints.NotNull;

@Embeddable
public record Token(@NotNull EncryptedToken value, @NotNull Instant expiresAtUtc) {
  public Token {
    Objects.requireNonNull(expiresAtUtc, "Token expiration time cannot be null");
    Objects.requireNonNull(value, "Token value must be provided");
  }

  public static CustomResult<Token, DomainError> create(EncryptedToken encryptedToken, Instant expiresAtUtc) {
    if (expiresAtUtc == null) {
      return CustomResult.failure(TokenErrors.TOKEN_EXPIRATION_TIME_CANNOT_BE_NULL);
    }
    return CustomResult.success(new Token(encryptedToken, expiresAtUtc));
  }
}
