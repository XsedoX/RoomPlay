package com.xsedox.roomplay_backend.domain.shared.result;

public record Success<TData, TError>(TData data) implements CustomResult<TData, TError> {
}
