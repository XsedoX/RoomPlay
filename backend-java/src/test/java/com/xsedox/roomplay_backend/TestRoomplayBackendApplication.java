package com.xsedox.roomplay_backend;

import org.springframework.boot.SpringApplication;

public class TestRoomplayBackendApplication {

	public static void main(String[] args) {
		SpringApplication.from(RoomplayBackendApplication::main).with(TestcontainersConfiguration.class).run(args);
	}

}
