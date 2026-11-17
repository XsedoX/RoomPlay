package com.xsedox.roomplay_backend.application.user.get_user;

import com.xsedox.roomplay_backend.application.common.Handler;
import com.xsedox.roomplay_backend.application.contracts.UserRepository;
import com.xsedox.roomplay_backend.application.errors.ApplicationError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;

/**
 * GetUserQueryHandler
 */
@Handler
public final class GetUserQueryHandler {
  private final UserRepository userRepository;

  public GetUserQueryHandler(UserRepository userRepository) {
    this.userRepository = userRepository;
  }

  public CustomResult<GetUserQueryResponse, ApplicationError> handle() {
    var user = userRepository.findById(query.userId());
  }
}
