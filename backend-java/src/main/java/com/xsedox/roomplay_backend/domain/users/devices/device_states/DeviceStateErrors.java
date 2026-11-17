package com.xsedox.roomplay_backend.domain.users.devices.device_states;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * DeviceStateErrors
 */
public final class DeviceStateErrors {

  public static final DomainError INVALID_DEVICE_STATE = new DomainError("DeviceStateErrors.INVALID_DEVICE_STATE",
      "The provided device state is invalid.");

}
