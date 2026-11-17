package com.xsedox.roomplay_backend.domain.users.devices.device_types;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * DeviceTypeErrors
 */
public final class DeviceTypeErrors {
  public static final DomainError INVALID_DEVICE_TYPE = new DomainError("DeviceTypeErrors.INVALID_DEVICE_TYPE",
      "The provided device type is invalid.");

}
