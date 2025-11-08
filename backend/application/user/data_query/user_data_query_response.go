package data_query

import "github.com/google/uuid"

type UserQueryResponse struct {
	Name    string     `json:"name"`
	Surname string     `json:"surname"`
	RoomId  *uuid.UUID `json:"roomId" example:"null" extensions:"x-nullable"`
	Role    *string    `json:"role"  example:"null" extensions:"x-nullable"`
}
