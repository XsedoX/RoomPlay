package com.xsedox.roomplay_backend.domain.users.devices;

import java.time.Instant;

import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import com.sun.istack.NotNull;
import com.xsedox.roomplay_backend.domain.shared.AbstractEntity;
import com.xsedox.roomplay_backend.domain.shared.DomainError;
import com.xsedox.roomplay_backend.domain.shared.result.CustomResult;
import com.xsedox.roomplay_backend.domain.users.devices.device_states.DeviceState;
import com.xsedox.roomplay_backend.domain.users.devices.device_types.DeviceType;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Table;

@Entity
@Table(name = "devices")
public class Device extends AbstractEntity<DeviceId> {
  protected Device() {
  }

  private Device(DeviceId id, String friendlyName, DeviceType deviceType, boolean isHost, DeviceState deviceState,
      Instant lastLoggedInUtc) {
    this.friendlyName = friendlyName;
    this.deviceType = deviceType;
    this.isHost = isHost;
    this.deviceState = deviceState;
    super(id);
  }

  @NotNull
  @Column(length = 30)
  private String friendlyName;
  @NotNull
  @JdbcTypeCode(SqlTypes.NAMED_ENUM)
  private DeviceType deviceType;
  private boolean isHost;
  @NotNull
  @JdbcTypeCode(SqlTypes.NAMED_ENUM)
  private DeviceState deviceState;
  @NotNull
  private Instant lastLoggedInUtc;

  public String getFriendlyName() {
    return friendlyName;
  }

  public DeviceType getDeviceType() {
    return deviceType;
  }

  public boolean isHost() {
    return isHost;
  }

  public DeviceState getDeviceState() {
    return deviceState;
  }

  public Instant getLastLoggedInUtc() {
    return lastLoggedInUtc;
  }

  public static CustomResult<Device, DomainError> create(DeviceType deviceType, Instant lastLoggedInUtc) {
    String frinedlyName;
    switch (deviceType) {
      case MOBILE:
        frinedlyName = "My lovely mobile device";
        break;
      case DESKTOP:
        frinedlyName = "My lovely computer";
        break;
      default:
        throw new IllegalArgumentException("Unknown device type: " + deviceType);
    }

    return CustomResult
        .success(new Device(DeviceId.create(), frinedlyName, deviceType, false, DeviceState.ONLINE, lastLoggedInUtc));
  }
}
