package com.xsedox.roomplay_backend.domain.users.devices.device_types;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

/**
 * DeviceType
 */
public enum DeviceType {
  MOBILE("mobile"),
  DESKTOP("desktop");

  private final String value;

  public String getValue() {
    return value;
  }

  DeviceType(String value) {
    this.value = value;
  }

  public static CustomResult<DeviceType, DomainError> fromValue(String raw) {
    for (DeviceType dt : values()) {
      if (dt.value.equals(raw))
        return CustomResult.success(dt);
    }
    return CustomResult.failure(DeviceTypeErrors.INVALID_DEVICE_TYPE);
  }
}
