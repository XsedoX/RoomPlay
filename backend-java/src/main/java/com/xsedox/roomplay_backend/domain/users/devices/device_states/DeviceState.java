package com.xsedox.roomplay_backend.domain.users.devices.device_states;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

public enum DeviceState {
  ONLINE("online"),
  OFFLINE("offline");

  private final String value;

  public String getValue() {
    return value;
  }

  DeviceState(String value) {
    this.value = value;
  }

  public static CustomResult<DeviceState, DomainError> fromValue(String raw) {
    for (DeviceState ds : values()) {
      if (ds.value.equals(raw))
        return CustomResult.success(ds);
    }
    return CustomResult.failure(DeviceStateErrors.INVALID_DEVICE_STATE);
  }
}
