package com.xsedox.roomplay_backend.domain.users.full_names;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

import jakarta.persistence.Embeddable;
import jakarta.validation.constraints.NotNull;

@Embeddable
public record FullName(
    @NotNull String name,
    @NotNull String surname) {
  public FullName {
    if (name == null || name.isBlank()) {
      throw new IllegalArgumentException("Name cannot be null or blank");
    }
    if (surname == null || surname.isBlank()) {
      throw new IllegalArgumentException("Surname cannot be null or blank");
    }
  }

  public static CustomResult<FullName, DomainError> create(String name, String surname) {
    if (name == null || name.isBlank()) {
      return CustomResult.failure(FullNameErrors.NAME_CANNOT_BE_NULL_OR_BLANK);
    }
    if (surname == null || surname.isBlank()) {
      return CustomResult.failure(FullNameErrors.SURNAME_CANNOT_BE_NULL_OR_BLANK);
    }
    return CustomResult.success(new FullName(name, surname));
  }
}
