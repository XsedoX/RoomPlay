package com.xsedox.roomplay_backend.infrastructure.persistance;

import org.springframework.data.jpa.repository.JpaRepository;

import com.xsedox.roomplay_backend.application.contracts.ExternalCredentialsRepository;
import com.xsedox.roomplay_backend.domain.external_credentials.ExternalCredentials;
import com.xsedox.roomplay_backend.domain.users.UserId;

/**
 * ConcreteExternalCredentialsRepository
 */
public interface JpaExternalCredentialsRepository
    extends ExternalCredentialsRepository, JpaRepository<ExternalCredentials, UserId> {
}
