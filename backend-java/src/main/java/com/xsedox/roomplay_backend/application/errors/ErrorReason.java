package com.xsedox.roomplay_backend.application.errors;

/**
 * ErrorReason
 */
public enum ErrorReason {
  UNKNOWN("Unknown error");

  private final String reason;

  public String getReason() {
    return reason;
  }

  ErrorReason(String reason) {
    this.reason = reason;
  }

  public static ErrorReason fromString(String reason) {
    for (ErrorReason errorReason : ErrorReason.values()) {
      if (errorReason.getReason().equals(reason)) {
        return errorReason;
      }
    }
    return UNKNOWN;
  }

}
