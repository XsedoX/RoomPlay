package com.xsedox.roomplay_backend.domain.rooms.qr_codes;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

import jakarta.persistence.Embeddable;
import jakarta.validation.constraints.NotNull;

/**
 * EncryptedQrCode
 */
@Embeddable
public record EncryptedQrCode(@NotNull String value) {
  public EncryptedQrCode {
    if (value == null || value.isBlank()) {
      throw new IllegalArgumentException(EncryptedQrCodeErrors.QR_CODE_STRING_BLANK_OR_NULL.message());
    }
  }

  public static CustomResult<EncryptedQrCode, DomainError> create(String value) {
    if (value == null || value.isBlank()) {
      return CustomResult.failure(EncryptedQrCodeErrors.QR_CODE_STRING_BLANK_OR_NULL);
    }
    return CustomResult.success(new EncryptedQrCode(value));
  }
}
