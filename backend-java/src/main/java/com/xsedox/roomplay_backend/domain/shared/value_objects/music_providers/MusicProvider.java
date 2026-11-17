package com.xsedox.roomplay_backend.domain.shared.value_objects.music_providers;

import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

public enum MusicProvider {
  SPOTIFY("spotify"),
  YOUTUBE("youtube");

  private final String value;

  public String getValue() {
    return value;
  }

  MusicProvider(String value) {
    this.value = value;
  }

  public static CustomResult<MusicProvider, DomainError> fromValue(String raw) {
    for (MusicProvider mp : values()) {
      if (mp.value.equals(raw))
        return CustomResult.success(mp);
    }
    return CustomResult.failure(MusicProviderErrors.MUSIC_PROVIDER_CANNOT_BE_NULL);
  }
}
