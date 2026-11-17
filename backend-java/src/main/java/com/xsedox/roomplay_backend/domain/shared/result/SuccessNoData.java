package com.xsedox.roomplay_backend.domain.shared.result;

public record SuccessNoData<TData, TError>() implements CustomResult<TData, TError> {
}
