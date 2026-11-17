package com.xsedox.roomplay_backend.domain.shared.result;

public sealed interface CustomResult<TData, TError> permits Success, Failure, SuccessNoData {

  public static <TData, TError> CustomResult<TData, TError> success(TData data) {
    return new Success<>(data);
  }

  public static <TData, TError> CustomResult<TData, TError> successNoData() {
    return new SuccessNoData<>();
  }

  public static <TData, TError> CustomResult<TData, TError> failure(TError error) {
    return new Failure<>(error);
  }
}
