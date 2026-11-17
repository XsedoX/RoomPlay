package com.xsedox.roomplay_backend.domain.shared.value_objects.music_providers;

import com.xsedox.roomplay_backend.domain.shared.DomainError;

/**
 * MusicProviderErros
 */
public final class MusicProviderErrors {
  public static final DomainError MUSIC_PROVIDER_CANNOT_BE_NULL = new DomainError(
      "MusicProviderErros.MUSIC_PROVIDER_CANNOT_BE_NULL",
      "Music provider cannot be null");

}
