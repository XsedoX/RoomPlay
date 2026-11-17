package com.xsedox.roomplay_backend.application.contracts;

import java.util.Optional;

import com.xsedox.roomplay_backend.domain.external_credentials.ExternalCredentials;
import com.xsedox.roomplay_backend.domain.users.UserId;

/**
 * JpaExternalCredentialsRepository
 */
public interface ExternalCredentialsRepository {
  Optional<ExternalCredentials> findById(UserId id);
}
