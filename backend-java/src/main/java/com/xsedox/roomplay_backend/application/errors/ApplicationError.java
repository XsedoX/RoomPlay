package com.xsedox.roomplay_backend.application.errors;

/**
 * ApplicationError
 */
public record ApplicationError(String code, String message, ErrorReason reason) {
}
