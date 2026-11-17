package com.xsedox.roomplay_backend.domain.users;

import java.time.Instant;
import java.util.Collections;
import java.util.HashSet;
import java.util.Optional;
import java.util.Set;

import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;
import org.jspecify.annotations.Nullable;

import com.xsedox.roomplay_backend.domain.rooms.RoomId;
import com.xsedox.roomplay_backend.domain.shared.AbstractAggregateRoot;
import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;
import com.xsedox.roomplay_backend.domain.shared.result.Failure;
import com.xsedox.roomplay_backend.domain.shared.result.Success;
import com.xsedox.roomplay_backend.domain.users.devices.Device;
import com.xsedox.roomplay_backend.domain.users.devices.device_types.DeviceType;
import com.xsedox.roomplay_backend.domain.users.full_names.FullName;
import com.xsedox.roomplay_backend.domain.users.user_roles.UserRole;

import jakarta.persistence.AttributeOverride;
import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Embedded;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.OneToMany;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotNull;

@Entity
@Table(name = "users")
public class User extends AbstractAggregateRoot<UserId> {
  @Embedded
  private FullName fullName;
  @Nullable
  @JdbcTypeCode(SqlTypes.NAMED_ENUM)
  private UserRole userRole;
  @Embedded
  @Nullable
  @AttributeOverride(name = "value", column = @Column(name = "room_id"))
  private RoomId roomId;
  @NotNull
  @OneToMany(cascade = CascadeType.ALL, orphanRemoval = true, fetch = FetchType.EAGER)
  @JoinColumn(name = "user_id")
  private Set<Device> devices = new HashSet<>();
  @Nullable
  private Instant boostUsedAtUtc;

  public FullName getFullName() {
    return fullName;
  }

  public Optional<UserRole> getUserRole() {
    return Optional.ofNullable(userRole);
  }

  public Optional<RoomId> getRoomId() {
    return Optional.ofNullable(roomId);
  }

  public Set<Device> getDevices() {
    return Collections.unmodifiableSet(devices);
  }

  public Optional<Instant> getBoostUsedAtUtc() {
    return Optional.ofNullable(boostUsedAtUtc);
  }

  protected User() {
  }

  private User(UserId id, FullName fullName, UserRole userRole, RoomId roomId, Set<Device> devices,
      Instant boostUsedAtUtc) {
    super(id);
    this.fullName = fullName;
    this.userRole = userRole;
    this.roomId = roomId;
    this.devices = devices;
    this.boostUsedAtUtc = boostUsedAtUtc;
  }

  public static CustomResult<User, DomainError> create(FullName fullName, DeviceType deviceType,
      Instant deviceLastLoggedInUtc) {
    CustomResult<Device, DomainError> deviceResult = Device.create(deviceType, deviceLastLoggedInUtc);
    Device device;
    switch (deviceResult) {
      case Success(var createdDevice):
        device = createdDevice;
        break;
      case Failure(var error):
        return CustomResult.failure(error);
      default:
        throw new IllegalStateException("Unexpected value: " + deviceResult);
    }
    User user = new User(UserId.create(), fullName, null, null, Set.of(device), null);
    return CustomResult.success(user);
  }
}
