package com.xsedox.roomplay_backend.domain.external_credentials;

import java.time.Instant;

import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import com.sun.istack.NotNull;
import com.xsedox.roomplay_backend.domain.shared.AbstractAggregateRoot;
import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;
import com.xsedox.roomplay_backend.domain.shared.value_objects.music_providers.MusicProvider;
import com.xsedox.roomplay_backend.domain.shared.value_objects.tokens.Token;
import com.xsedox.roomplay_backend.domain.users.UserId;

import jakarta.persistence.AttributeOverride;
import jakarta.persistence.AttributeOverrides;
import jakarta.persistence.Column;
import jakarta.persistence.Embedded;
import jakarta.persistence.Entity;
import jakarta.persistence.Table;

@Entity
@Table(name = "external_credentials")
@AttributeOverride(name = "id.value", column = @Column(name = "user_id"))
public class ExternalCredentials extends AbstractAggregateRoot<UserId> {
  @NotNull
  @Column(unique = true)
  private String externalId;
  @Embedded
  @AttributeOverrides({
      @AttributeOverride(name = "token.value", column = @Column(name = "access_token")),
      @AttributeOverride(name = "token.expiresAtUtc", column = @Column(name = "access_token_expires_at_utc"))
  })
  private Token accessToken;
  @NotNull
  private String refreshToken;
  @NotNull
  @JdbcTypeCode(SqlTypes.NAMED_ENUM)
  private MusicProvider musicProvider;
  @NotNull
  private Instant accessTokenIssuedAtUtc;

  protected ExternalCredentials() {
  }

  private ExternalCredentials(UserId id, String externalId, Token accessToken, String refreshToken,
      MusicProvider musicProvider, Instant accessTokenIssuedAtUtc) {
    super(id);
    this.externalId = externalId;
    this.accessToken = accessToken;
    this.refreshToken = refreshToken;
    this.musicProvider = musicProvider;
    this.accessTokenIssuedAtUtc = accessTokenIssuedAtUtc;
  }

  public static CustomResult<ExternalCredentials, DomainError> create(UserId id, String externalId, Token accessToken,
      String refreshToken, MusicProvider musicProvider, Instant accessTokenIssuedAtUtc) {
    if (externalId == null || externalId.isBlank()) {
      return CustomResult.failure(ExternalCredentialsErrors.EXTERNAL_ID_CANNOT_BE_NULL_OR_BLANK);
    }
    if (accessToken == null) {
      return CustomResult.failure(ExternalCredentialsErrors.ACCESS_TOKEN_CANNOT_BE_NULL);
    }
    if (refreshToken == null || refreshToken.isBlank()) {
      return CustomResult.failure(ExternalCredentialsErrors.REFRESH_TOKEN_CANNOT_BE_NULL_OR_BLANK);
    }
    if (musicProvider == null) {
      return CustomResult.failure(ExternalCredentialsErrors.MUSIC_PROVIDER_CANNOT_BE_NULL);
    }
    if (accessTokenIssuedAtUtc == null) {
      return CustomResult.failure(ExternalCredentialsErrors.ACCESS_TOKEN_ISSUED_AT_UTC_CANNOT_BE_NULL);
    }

    return CustomResult.success(new ExternalCredentials(id, externalId, accessToken, refreshToken, musicProvider,
        accessTokenIssuedAtUtc));
  }

  public String getExternalId() {
    return externalId;
  }

  public String getRefreshToken() {
    return refreshToken;
  }

  public MusicProvider getMusicProvider() {
    return musicProvider;
  }

  public Instant getAccessTokenIssuedAtUtc() {
    return accessTokenIssuedAtUtc;
  }

  public Token getAccessToken() {
    return accessToken;
  }
}
