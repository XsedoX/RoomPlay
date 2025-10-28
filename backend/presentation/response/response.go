package response

type Success struct {
	Message *string `json:"message"  example:"null" extensions:"x-nullable"`
	Data    any     `json:"data" swaggertype:"object"`
}

func Ok(data any) Success {
	return Success{
		Message: nil,
		Data:    data,
	}
}

type Error struct {
	Message string  `json:"message"`
	Data    *string `json:"data" example:"null" extensions:"x-nullable"`
}

func Failure(message string) Error {
	return Error{
		Message: message,
		Data:    nil,
	}
}
