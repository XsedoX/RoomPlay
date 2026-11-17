package com.xsedox.roomplay_backend.domain.users.user_roles;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

public enum UserRole {
  HOST("host"),
  MEMBER("member");

  private final String value;

  public String getValue() {
    return value;
  }

  UserRole(String value) {
    this.value = value;
  }

  public static CustomResult<UserRole, DomainError> fromValue(String raw) {
    for (UserRole ur : values()) {
      if (ur.value.equals(raw))
        return CustomResult.success(ur);
    }
    return CustomResult.failure(UserRoleErrors.INVALID_USER_ROLE);
  }
}
