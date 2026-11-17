package com.xsedox.roomplay_backend.application.contracts;

import java.util.Optional;

import com.xsedox.roomplay_backend.domain.users.User;
import com.xsedox.roomplay_backend.domain.users.UserId;

/**
 * UserRepository
 */
public interface UserRepository {
  Optional<User> findById(UserId id);
}
