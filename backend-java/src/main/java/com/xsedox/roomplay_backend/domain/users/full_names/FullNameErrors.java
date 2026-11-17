package com.xsedox.roomplay_backend.domain.users.full_names;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

public final class FullNameErrors {
  public static final DomainError NAME_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "FullNameErrors.NAME_CANNOT_BE_NULL_OR_BLANK", "Name must be provided");
  public static final DomainError SURNAME_CANNOT_BE_NULL_OR_BLANK = new DomainError(
      "FullNameErrors.SURNAME_CANNOT_BE_NULL_OR_BLANK",
      "Surname must be provided");
}
