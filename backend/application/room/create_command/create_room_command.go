package create_command

type CreateRoomCommand struct {
	RoomName           string `json:"roomName" validate:"required,gte=5,lte=30"`
	RoomPassword       string `json:"roomPassword" validate:"required,gte=10,lte=30,no_whitespace"`
	RepeatRoomPassword string `json:"repeatRoomPassword" validate:"required,eqcsfield=RoomPassword"`
}
