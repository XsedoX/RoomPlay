package com.xsedox.roomplay_backend.domain.shared.result;

import java.util.Objects;

public record Failure<TData, TError>(TError error) implements CustomResult<TData, TError> {
  public Failure {
    Objects.requireNonNull(error, "error must not be null");
  }
}
