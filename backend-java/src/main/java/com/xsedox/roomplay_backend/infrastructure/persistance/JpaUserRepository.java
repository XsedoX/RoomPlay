package com.xsedox.roomplay_backend.infrastructure.persistance;

import org.springframework.data.jpa.repository.JpaRepository;

import com.xsedox.roomplay_backend.application.contracts.UserRepository;
import com.xsedox.roomplay_backend.domain.users.User;
import com.xsedox.roomplay_backend.domain.users.UserId;

/**
 * JpaUserRepository
 */
public interface JpaUserRepository extends UserRepository, JpaRepository<User, UserId> {
}
