package join_room_password

type JoinRoomPasswordCommand struct {
	RoomName     string `json:"roomName" fname:"Room Name" validate:"required,gte=5,lte=30"`
	RoomPassword string `json:"roomPassword" fname:"Room Password" validate:"required,gte=10,lte=30,no_whitespace"`
}
