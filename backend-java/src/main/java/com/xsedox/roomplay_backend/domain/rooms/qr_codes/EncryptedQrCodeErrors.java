package com.xsedox.roomplay_backend.domain.rooms.qr_codes;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * EncryptedQrCodeErrors
 */
public final class EncryptedQrCodeErrors {
  public static DomainError QR_CODE_STRING_BLANK_OR_NULL = new DomainError(
      "EncryptedQrCodeErrors.QR_CODE_STRING_BLANK_OR_NULL", "The QR code must be provided");
}
