package com.xsedox.roomplay_backend.domain.users.user_roles;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * UserRoleErrors
 */
public final class UserRoleErrors {
  public static final DomainError INVALID_USER_ROLE = new DomainError("UserRoleErrors.INVALID_USER_ROLE",
      "The provided user role is invalid.");

}
