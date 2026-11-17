package com.xsedox.roomplay_backend.domain.external_credentials;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

public final class ExternalCredentialsErrors {
  public static final DomainError EXTERNAL_ID_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "ExternalCredentialsErrors.EXTERNAL_ID_CANNOT_BE_NULL_OR_BLANK",
      "External ID has to be provided and cannot be blank");
  public static final DomainError ACCESS_TOKEN_CANNOT_BE_NULL = new DomainError(
      "ExternalCredentialsErrors.ACCESS_TOKEN_CANNOT_BE_NULL", "Access token has to be provided");
  public static final DomainError REFRESH_TOKEN_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "ExternalCredentialsErrors.REFRESH_TOKEN_CANNOT_BE_NULL_OR_BLANK",
      "Refresh token has to be provided and cannot be blank");
  public static final DomainError MUSIC_PROVIDER_CANNOT_BE_NULL = new DomainError(
      "ExternalCredentialsErrors.MUSIC_PROVIDER_CANNOT_BE_NULL", "Music provider has to be provided");
  public static final DomainError ACCESS_TOKEN_ISSUED_AT_UTC_CANNOT_BE_NULL = new DomainError(
      "ExternalCredentialsErrors.ACCESS_TOKEN_ISSUED_AT_UTC_CANNOT_BE_NULL",
      "Access token issued at UTC has to be provided");
}
